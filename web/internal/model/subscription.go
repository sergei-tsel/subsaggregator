package model

import (
	_ "subsaggregator/docs"
	"subsaggregator/internal/utils"
)

// Subscription представляет запись о подписке
//
//	@modelId	sub
//
// swagger:model Subscription
type Subscription struct {
	// ИД записи о подписке
	// required: true
	// min: 1
	Id int `json:"id"`

	// Название сервиса, предоставляющего подписку
	// required: true
	// example: "Yandex Plus"
	ServiceName string `json:"service_name"`

	// Стоимость месячной подписки в рублях
	// required: true
	// min: 1
	Price int `json:"price"`

	// ИД пользователя
	// required: true
	// min: 1
	UserId string `json:"user_id"`

	// Дата начала подписки
	// required: true
	StartDate *utils.Date `json:"start_date"`

	// Дата окончания подписки
	// required: false
	EndDate *utils.Date `json:"end_date,omitempty"`
}
