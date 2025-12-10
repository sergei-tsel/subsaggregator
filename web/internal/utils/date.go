package utils

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Date переопределяет NullTime из database/sql
type Date struct {
	sql.NullTime
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`null`), nil
	}

	return json.Marshal(d.Time.Format("01-2006"))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Time = time.Time{}

		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsedTime, err := time.Parse("01-2006", s)

	if err != nil {
		return err
	}

	*d = Date{NullTime: sql.NullTime{
		Time:  parsedTime,
		Valid: true,
	}}

	return nil
}
