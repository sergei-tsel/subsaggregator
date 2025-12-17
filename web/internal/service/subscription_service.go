package service

import (
	_ "subsaggregator/docs"
	"subsaggregator/internal/model"
	"subsaggregator/internal/repository"
	"subsaggregator/internal/utils"
)

// CreateSubscriptionRequest Модель данных для создания записи о подписке
//
//	@modelId	create-sub-request
//	@required	ServiceName Price UserId StartDate EndDate
type CreateSubscriptionRequest struct {
	ServiceName string      `json:"service_name"`
	Price       int         `json:"price"`
	UserId      string      `json:"user_id"`
	StartDate   *utils.Date `json:"start_date" swaggertype:"string" example:"07-2025"`
	EndDate     *utils.Date `json:"end_date,omitempty" swaggertype:"string" example:"07-2025"`
}

// UpdateSubscriptionRequest Модель данных для изменения записи о подписке
//
//	@modelId	update-sub-request
//	@required	ServiceName Price UserId StartDate EndDate
type UpdateSubscriptionRequest struct {
	ServiceName string      `json:"service_name"`
	Price       int         `json:"price"`
	UserId      string      `json:"user_id"`
	StartDate   *utils.Date `json:"start_date" swaggertype:"string" example:"07-2025"`
	EndDate     *utils.Date `json:"end_date,omitempty" swaggertype:"string" example:"07-2025"`
}

// ListSubscriptionsRequest Модель данных для получения списка записей о подписках
//
//	@modelId	list-subs-request
//	@required	ServiceName UserId StartDate EndDate Offset Limit
type ListSubscriptionsRequest struct {
	ServiceName string     `json:"service_name,omitempty"`
	UserId      string     `json:"user_id,omitempty"`
	StartDate   utils.Date `json:"start_date" swaggertype:"string" example:"07-2025"`
	EndDate     utils.Date `json:"end_date,omitempty" swaggertype:"string" example:"07-2025"`
	Offset      int        `json:"offset"`
	Limit       int        `json:"limit" example:"10"`
}

// SumSubscriptionsPricesRequest Модель данных для получения суммарной стоимости подписок
//
//	@modelId	sum-subs-prices-request
//	@required	ServiceName UserId StartDate EndDate
type SumSubscriptionsPricesRequest struct {
	ServiceName string     `json:"service_name,omitempty"`
	UserId      string     `json:"user_id,omitempty"`
	StartDate   utils.Date `json:"start_date" swaggertype:"string" example:"07-2025"`
	EndDate     utils.Date `json:"end_date,omitempty" swaggertype:"string" example:"07-2025"`
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
		req.Offset,
		req.Limit,
	)

	if err != nil {
		return subs, err
	}

	return subs, nil
}

func SumSubscriptionsPrices(req SumSubscriptionsPricesRequest, subscriptionRepo repository.SubscriptionRepository) (*int, error) {
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

	sub.ServiceName = req.ServiceName
	sub.Price = req.Price
	sub.UserId = req.UserId
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
