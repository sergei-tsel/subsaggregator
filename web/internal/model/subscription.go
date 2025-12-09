package model

import (
	"subsaggregator/internal/utils"
)

type Subscription struct {
	Id          int         `json:"id"`                 // уникальный идентификатор
	ServiceName string      `json:"service_name"`       // название сервиса, предоставляющего подписку
	Price       int         `json:"price"`              // стоимость месячной подписки в рублях
	UserId      string      `json:"user_id"`            // уникальный идентификатор пользователя
	StartDate   *utils.Date `json:"start_date"`         // дата начала подписки
	EndDate     *utils.Date `json:"end_date,omitempty"` // дата окончания подписки
}
