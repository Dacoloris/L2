package dev11

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type Service struct {
	Cache *Cache
}

func NewService(cache *Cache) *Service {
	return &Service{
		Cache: cache,
	}
}

func (s *Service) CreateEvent(data []byte) error {
	var e Event
	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	s.Cache.Add(e)
	return nil
}

func (s *Service) UpdateEvent(data []byte) error {
	var e Event
	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	if !s.Cache.Contains(e) {
		return fmt.Errorf("event id %s not found", e.ID)
	}

	s.Cache.Update(e)
	return nil
}

func (s *Service) DeleteEvent(data []byte) error {
	var e Event
	err := json.Unmarshal(data, &e)
	if err != nil {
		return fmt.Errorf("event id %s not found", e.ID)
	}
	s.Cache.Delete(e)
	return nil
}

func (s *Service) GetEventsForDay(params url.Values) ([]byte, error) {
	date, err := time.Parse(dateLayout, params.Get("date"))
	if err != nil {
		return nil, err
	}

	events := s.Cache.GetByDate(date)
	jsonRes := NewEventsResponse(events)

	return jsonRes, nil
}

func (s *Service) GetEventsForWeek(params url.Values) ([]byte, error) {
	date, err := time.Parse(dateLayout, params.Get("date"))
	if err != nil {
		return nil, err
	}

	events := s.Cache.GetByWeek(date)
	jsonRes := NewEventsResponse(events)

	return jsonRes, nil
}

func (s *Service) GetEventsForMonth(params url.Values) ([]byte, error) {
	date, err := time.Parse(dateLayout, params.Get("date"))
	if err != nil {
		return nil, err
	}

	events := s.Cache.GetByMonth(date)
	jsonRes := NewEventsResponse(events)

	return jsonRes, nil
}
