package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"subsaggregator/internal/db"
	"subsaggregator/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type SubscriptionRepository interface {
	FindById(id int) (*model.Subscription, error)
	List(userId string, serviceName string) ([]model.Subscription, error)
	SumPrices(userId string, serviceName string) (*int, error)
	Create(entity *model.Subscription) error
	Update(entity *model.Subscription) error
	Delete(entity *model.Subscription) error
}

type SubscriptionRepo struct{}

func (repo *SubscriptionRepo) FindById(id int) (*model.Subscription, error) {
	subCache, err := getSubscriptionCache(id)

	if subCache != nil {
		return subCache, nil
	}

	query := `
		SELECT *
		FROM subscriptions
		WHERE id = $1;
	`

	row, err := db.Postgres.Query(query, id)
	defer row.Close()

	if !row.Next() {
		return nil, fmt.Errorf("subsription not found: %w", err)
	}

	var sub model.Subscription

	err = row.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate, &sub.EndDate)

	if err != nil {
		return nil, fmt.Errorf("failing to read data from database: %w", err)
	}

	setSubscriptionCache(&sub)

	return &sub, nil
}

func (repo *SubscriptionRepo) List(userId string, serviceName string) ([]model.Subscription, error) {
	query := `
    	SELECT subscriptions.*
		FROM subscriptions
    	WHERE ($1 IS NULL OR subscriptions.user_id = $1) AND ($2 IS NULL OR subscriptions.service_name = $2);
	`

	rows, err := db.Postgres.Query(query, userId, serviceName)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("subscriptions not found: %w", err)
	}

	var subs []model.Subscription

	for rows.Next() {
		var sub model.Subscription

		err = rows.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate, &sub.EndDate)

		if err != nil {
			return nil, fmt.Errorf("failing to read data from database: %w", err)
		}

		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failing to read data from database: %w", err)
	}

	return subs, nil
}

func (repo *SubscriptionRepo) SumPrices(userId string, serviceName string) (*int, error) {
	query := `
    	SELECT SUM(subscriptions.price) OVER () AS total_price
		FROM subscriptions
    	WHERE ($1 IS NULL OR subscriptions.user_id = $1) AND ($2 IS NULL OR subscriptions.service_name = $2);
	`

	row := db.Postgres.QueryRow(query, userId, serviceName)

	var sumPrice int

	err := row.Scan(&sumPrice)

	if err != nil {
		return nil, fmt.Errorf("failed to sum subscriptions prices: %w", err)
	}

	return &sumPrice, nil
}

func (repo *SubscriptionRepo) Create(entity *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Postgres.Query(
		query,
		entity.ServiceName,
		entity.Price,
		entity.UserId,
		entity.StartDate,
		entity.EndDate,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	setSubscriptionCache(entity)

	return nil
}

func (repo *SubscriptionRepo) Update(entity *model.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET service_name = $2, price = $3, user_id = $4, start_date = $5, end_date = $6
		WHERE id = $1
	`

	_, err := db.Postgres.Query(
		query,
		entity.Id,
		entity.ServiceName,
		entity.Price,
		entity.UserId,
		entity.StartDate,
		entity.EndDate,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	setSubscriptionCache(entity)

	return nil
}

func (repo *SubscriptionRepo) Delete(entity *model.Subscription) error {
	query := `
		DELETE 
		FROM subscriptions 
		WHERE id = $1;
	`

	_, err := db.Postgres.Query(query, entity.Id)

	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	deleteSubscriptionCache(entity)

	return nil
}

func setSubscriptionCache(sub *model.Subscription) {
	jsonBytes, _ := json.Marshal(sub)

	db.Redis.Set(
		context.Background(),
		fmt.Sprintf("sub:%d", sub.Id),
		jsonBytes,
		3*time.Minute,
	)
}

func getSubscriptionCache(subId int) (*model.Subscription, error) {
	result, err := db.Redis.Get(
		context.Background(),
		fmt.Sprintf("sub:%d", subId),
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, err
	}

	var sub model.Subscription

	err = json.Unmarshal([]byte(result), &sub)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func deleteSubscriptionCache(sub *model.Subscription) {
	db.Redis.Del(
		context.Background(),
		fmt.Sprintf("sub:%d", sub.Id),
	)
}
