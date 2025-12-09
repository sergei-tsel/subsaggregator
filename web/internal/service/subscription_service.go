package service

import (
	"subsaggregator/internal/model"
	"subsaggregator/internal/repository"
	"subsaggregator/internal/utils"
)

type CreateSubscriptionRequest struct {
	ServiceName string      `json:"service_name"`
	Price       int         `json:"price"`
	UserId      string      `json:"user_id"`
	StartDate   *utils.Date `json:"start_date"`
	EndDate     *utils.Date `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	Price     int         `json:"price"`
	StartDate *utils.Date `json:"start_date"`
	EndDate   *utils.Date `json:"end_date,omitempty"`
}

type ListSubscriptionsRequest struct {
	ServiceName string     `json:"service_name,omitempty"`
	UserId      string     `json:"user_id,omitempty"`
	StartDate   utils.Date `json:"start_date"`
	EndDate     utils.Date `json:"end_date,omitempty"`
}

type SumSubscriptionsRequest struct {
	ServiceName string     `json:"service_name,omitempty"`
	UserId      string     `json:"user_id,omitempty"`
	StartDate   utils.Date `json:"start_date"`
	EndDate     utils.Date `json:"end_date,omitempty"`
}

func GetOneSubscription(subscriptionRepo repository.SubscriptionRepository, subsId int) (*model.Subscription, error) {
	var sub *model.Subscription

	sub, err := subscriptionRepo.FindById(subsId)

	if err != nil {
		return sub, err
	}

	return sub, nil
}

func ListSubscriptions(req ListSubscriptionsRequest, subscriptionRepo repository.SubscriptionRepository) ([]model.Subscription, error) {
	subs, err := subscriptionRepo.List(
		req.UserId,
		req.ServiceName,
		req.StartDate,
		req.EndDate,
	)

	if err != nil {
		return subs, err
	}

	return subs, nil
}

func SumSubscriptionsPrices(req SumSubscriptionsRequest, subscriptionRepo repository.SubscriptionRepository) (*int, error) {
	sum, err := subscriptionRepo.SumPrices(
		req.UserId,
		req.ServiceName,
		req.StartDate,
		req.EndDate,
	)

	if err != nil {
		return sum, err
	}

	return sum, nil
}

func CreateSubscription(req CreateSubscriptionRequest, repo repository.SubscriptionRepository) (*model.Subscription, error) {
	sub := &model.Subscription{}
	sub.ServiceName = req.ServiceName
	sub.Price = req.Price
	sub.UserId = req.UserId
	sub.StartDate = req.StartDate
	sub.EndDate = req.EndDate

	err := repo.Create(sub)

	if err != nil {
		return sub, err
	}

	return sub, nil
}

func UpdateSubscription(req UpdateSubscriptionRequest, repo repository.SubscriptionRepository, subsId int) (*model.Subscription, error) {
	sub, err := GetOneSubscription(repo, subsId)

	if sub == nil {
		return nil, err
	}

	sub.Price = req.Price
	sub.StartDate = req.StartDate
	sub.EndDate = req.EndDate

	err = repo.Update(sub)

	if err != nil {
		return sub, err
	}

	return sub, nil
}

func DeleteSubscription(repo repository.SubscriptionRepository, subId int) error {
	sub, err := GetOneSubscription(repo, subId)

	if sub == nil {
		return err
	}

	err = repo.Delete(sub)

	if err != nil {
		return err
	}

	return nil
}
