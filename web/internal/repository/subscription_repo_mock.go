package repository

import (
	"fmt"
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

func (repo SubscriptionRepoMock) List(
	userId string,
	serviceName string,
	maxStartDate utils.Date,
	minEndDate utils.Date,
	offset int,
	limit int,
) ([]model.Subscription, error) {
	subs := []model.Subscription{}
	count := 0

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

		if count < offset {
			count++
			continue
		}

		subs = append(subs, *sub)

		if len(subs) == limit {
			return subs, nil
		}
	}

	return subs, nil
}

func (repo SubscriptionRepoMock) SumPrices(userId string, serviceName string, maxStartDate utils.Date, minEndDate utils.Date) (*int, error) {
	var sumPrice int

	uniquePrices := make(map[string]bool)

	for _, sub := range repo.Subscriptions {
		if maxStartDate.NullTime.Time.After(sub.EndDate.Time) ||
			minEndDate.NullTime.Time.Before(sub.StartDate.Time) {
			continue
		}

		if userId != "" && sub.UserId != userId {
			continue
		}

		if serviceName != "" && sub.ServiceName != serviceName {
			continue
		}

		for currentMonth := sub.StartDate.NullTime.Time; !currentMonth.After(sub.EndDate.NullTime.Time); currentMonth = currentMonth.AddDate(0, 1, 0) {
			key := fmt.Sprintf("%s:%s:%s", sub.UserId, sub.ServiceName, currentMonth.Format("01-2006"))

			if _, exists := uniquePrices[key]; !exists {
				sumPrice += sub.Price
				uniquePrices[key] = true
			}
		}
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
