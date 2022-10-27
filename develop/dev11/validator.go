package dev11

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateEventCreate(data []byte) error {
	var e Event
	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) ValidateEventID(data []byte) error {
	var e Event
	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}
	if e.ID == "" {
		return fmt.Errorf("event id is not provided")
	}

	return nil
}

func (v *Validator) ValidateDate(params url.Values) error {
	if !params.Has("date") {
		return fmt.Errorf("no date provided")
	}

	_, err := time.Parse(dateLayout, params.Get("date"))
	if err != nil {
		return err
	}

	return nil
}
