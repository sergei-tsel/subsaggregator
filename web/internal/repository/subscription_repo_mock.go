package repository

import (
	"subsaggregator/internal/model"
	"subsaggregator/internal/utils"
)

type SubscriptionRepoMock struct {
	Subscriptions map[int]*model.Subscription
	Count         int
}

func (repo SubscriptionRepoMock) FindById(id int) (*model.Subscription, error) {
	if repo.Subscriptions[id] != nil {
		return repo.Subscriptions[id], nil
	}

	return nil, nil
}

func (repo SubscriptionRepoMock) List(userId string, serviceName string, maxStartDate utils.Date, minEndDate utils.Date) ([]model.Subscription, error) {
	subs := []model.Subscription{}

	for _, sub := range repo.Subscriptions {
		if maxStartDate.NullTime.Time.After(sub.EndDate.Time) || minEndDate.NullTime.Time.Before(sub.StartDate.Time) {
			continue
		}

		if userId != "" && sub.UserId != userId {
			continue
		}

		if serviceName != "" && sub.ServiceName != serviceName {
			continue
		}

		subs = append(subs, *sub)
	}

	return subs, nil
}

func (repo SubscriptionRepoMock) SumPrices(userId string, serviceName string, maxStartDate utils.Date, minEndDate utils.Date) (*int, error) {
	var sumPrice int

	for _, sub := range repo.Subscriptions {
		if maxStartDate.NullTime.Time.After(sub.EndDate.Time) || minEndDate.NullTime.Time.Before(sub.StartDate.Time) {
			continue
		}

		if userId != "" && sub.UserId != userId {
			continue
		}

		if serviceName != "" && sub.ServiceName != serviceName {
			continue
		}

		sumPrice = sumPrice + sub.Price
	}

	return &sumPrice, nil
}

func (repo SubscriptionRepoMock) Create(entity *model.Subscription) error {
	repo.Count++

	entity.Id = repo.Count

	repo.Subscriptions[entity.Id] = entity

	return nil
}

func (repo SubscriptionRepoMock) Update(entity *model.Subscription) error {
	repo.Subscriptions[entity.Id].Price = entity.Price
	repo.Subscriptions[entity.Id].StartDate = entity.StartDate
	repo.Subscriptions[entity.Id].EndDate = entity.EndDate

	return nil
}

func (repo SubscriptionRepoMock) Delete(entity *model.Subscription) error {
	delete(repo.Subscriptions, entity.Id)

	return nil
}
