package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"subsaggregator/internal/db"
	"subsaggregator/internal/model"
	"subsaggregator/internal/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type SubscriptionRepository interface {
	FindById(id int) (*model.Subscription, error)
	List(userId string, serviceName string, maxStartDate utils.Date, minEndDate utils.Date, offset int, limit int) ([]model.Subscription, error)
	SumPrices(userId string, serviceName string, maxStartDate utils.Date, minEndDate utils.Date) (*int, error)
	Create(entity *model.Subscription) error
	Update(entity *model.Subscription) error
	Delete(entity *model.Subscription) error
}

type SubscriptionRepo struct{}

func (repo *SubscriptionRepo) FindById(id int) (*model.Subscription, error) {
	subCache, err := getSubscriptionCache(id)
	if subCache != nil {
		slog.Info(fmt.Sprintf("Получение кешированной записи о подписке. ИД: %d", id))

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
		slog.Error(fmt.Errorf("запись о подписке не найдена: %w", err).Error())

		return nil, fmt.Errorf("subsription not found: %w", err)
	}

	var sub model.Subscription

	err = row.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate, &sub.EndDate)

	if err != nil {
		slog.Error(fmt.Errorf("запись о подписке невозможно прочитать: %w", err).Error())

		return nil, fmt.Errorf("failing to read data from database: %w", err)
	}

	setSubscriptionCache(&sub)

	slog.Info(fmt.Sprintf("Получение записи о подписке. ИД: %d", id))

	return &sub, nil
}

func (repo *SubscriptionRepo) List(
	userId string,
	serviceName string,
	maxStartDate utils.Date,
	minEndDate utils.Date,
	offset int,
	limit int,
) ([]model.Subscription, error) {
	query := `
    	SELECT subscriptions.*
		FROM subscriptions
    	WHERE ($1::TEXT IS NULL OR subscriptions.user_id = $1) 
    	  AND ($2::TEXT IS NULL OR subscriptions.service_name = $2)
    	  AND (CASE 
       			WHEN $3 <> '0001-01-01'::DATE THEN subscriptions.start_date <= $3
       			ELSE TRUE
     		END)
    	  AND (CASE 
       			WHEN $4 <> '0001-01-01'::DATE THEN subscriptions.end_date <= $4
       			ELSE TRUE
     		END)
    	OFFSET $5
    	LIMIT $6;
	`

	rows, err := db.Postgres.Query(query, getFilter(userId), getFilter(serviceName), maxStartDate, minEndDate, offset, limit)
	defer rows.Close()

	if err != nil {
		slog.Error(fmt.Errorf("записи о подписках не найдены: %w", err).Error())

		return nil, fmt.Errorf("subscriptions not found: %w", err)
	}

	var subs []model.Subscription

	for rows.Next() {
		var sub model.Subscription

		err = rows.Scan(&sub.Id, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate, &sub.EndDate)

		if err != nil {
			slog.Error(fmt.Errorf("запись о подписке невозможно прочитать: %w", err).Error())

			return nil, fmt.Errorf("failing to read data from database: %w", err)
		}

		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		slog.Error(fmt.Errorf("записи о подписке невозможно прочитать: %w", err).Error())

		return nil, fmt.Errorf("failing to read data from database: %w", err)
	}

	slog.Info(fmt.Sprintf(
		"Получение записей о подписках c %s до %s. ИД пользователя: %s. Название сервиса: %s",
		minEndDate.Time.Format("01-2006"),
		maxStartDate.Time.Format("01-2006"),
		userId,
		serviceName,
	))

	return subs, nil
}

func (repo *SubscriptionRepo) SumPrices(userId string, serviceName string, maxStartDate utils.Date, minEndDate utils.Date) (*int, error) {
	query := `
    	SELECT SUM(subscriptions.price) OVER () AS total_price
		FROM subscriptions
    	WHERE ($1::TEXT IS NULL OR subscriptions.user_id = $1) 
    	  AND ($2::TEXT IS NULL OR subscriptions.service_name = $2)
    	  AND (CASE 
       			WHEN $3 <> '0001-01-01'::DATE THEN subscriptions.start_date <= $3
       			ELSE TRUE
     		END)
    	  AND (CASE 
    	      	WHEN $4 <> '0001-01-01'::DATE THEN subscriptions.end_date <= $4
       			ELSE TRUE
     		END);
	`

	row := db.Postgres.QueryRow(query, getFilter(userId), getFilter(serviceName), maxStartDate, minEndDate)

	var sumPrice int

	err := row.Scan(&sumPrice)

	if err != nil {
		slog.Error(fmt.Errorf("суммарная стоимость подписок не получена: %w", err).Error())

		return nil, fmt.Errorf("failed to sum subscriptions prices: %w", err)
	}

	slog.Info(fmt.Sprintf(
		"Получение суммарной стоимости подписок c %s до %s. ИД пользователя: %s. Название сервиса: %s",
		minEndDate.Time.Format("01-2006"),
		maxStartDate.Time.Format("01-2006"),
		userId,
		serviceName,
	))

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
		slog.Error(fmt.Errorf("запись о подписке не создана: %w", err).Error())

		return fmt.Errorf("failed to create subscription: %w", err)
	}

	setSubscriptionCache(entity)

	slog.Info(fmt.Sprintf("Создание записи о подписке. Название сервиса: %s. Стоимость: %d. ИД пользователя: %s. Дата начала: %s. Дата окончания: %s",
		entity.ServiceName,
		entity.Price,
		entity.UserId,
		entity.StartDate,
		entity.EndDate,
	))

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
		slog.Error(fmt.Errorf("запись о подписке не изменена: %w", err).Error())

		return fmt.Errorf("failed to update subscription: %w", err)
	}

	setSubscriptionCache(entity)

	slog.Info(fmt.Sprintf("Изменение записи о подписке. Название сервиса: %s. Стоимость: %d. ИД пользователя: %s. Дата начала: %s. Дата окончания: %s",
		entity.ServiceName,
		entity.Price,
		entity.UserId,
		entity.StartDate,
		entity.EndDate,
	))

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
		slog.Error(fmt.Errorf("запись о подписке не удалена: %w", err).Error())

		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	deleteSubscriptionCache(entity)

	slog.Info(fmt.Sprintf("Удаление записи о подписке. Название сервиса: %s. Стоимость: %d. ИД пользователя: %s. Дата начада: %s. Дата окончания: %s",
		entity.ServiceName,
		entity.Price,
		entity.UserId,
		entity.StartDate,
		entity.EndDate,
	))

	return nil
}

func getFilter(value string) sql.NullString {
	var filter sql.NullString
	if value == "" {
		filter.Valid = false
	} else {
		filter.String = value
		filter.Valid = true
	}

	return filter
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
