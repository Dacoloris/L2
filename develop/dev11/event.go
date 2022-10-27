package dev11

import (
	"bytes"
	"fmt"
	"time"
)

type Event struct {
	ID      string `json:"id"`
	Date    Date   `json:"date"`
	Content string `json:"content"`
}

type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d *Date) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	if len(data) == 0 {
		return fmt.Errorf("invalid date format: '%s'", string(data))
	}

	date, err := time.Parse(dateLayout, string(data))
	if err != nil {
		return fmt.Errorf("invalid date format: '%s'", string(data))
	}
	d.Time = date
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	dateStr := d.Format(dateLayout)
	stamp := fmt.Sprintf(`%q`, dateStr)
	return []byte(stamp), nil
}
