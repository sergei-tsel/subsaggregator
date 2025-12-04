package model

import (
	"encoding/json"
	"time"
)

type Date struct {
	time.Time
}

func (date Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(date.Format("07-2025"))
}

type Subscription struct {
	Id          int    // уникальный идентификатор
	ServiceName string // название сервиса, предоставляющего подписку
	Price       int    // стоимость месячной подписки в рублях
	UserId      string // уникальный идентификатор пользователя
	StartDate   Date   // дата начала подписки
	EndDate     Date   // дата окончания подписки
}
