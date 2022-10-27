package dev11

import (
	"strconv"
	"time"
)

type Cache struct {
	ID     int
	Events map[string]Event
	Tree   map[int]map[time.Month]map[int][]string
}

func NewCache() *Cache {
	return &Cache{
		ID:     0,
		Events: map[string]Event{},
		Tree:   map[int]map[time.Month]map[int][]string{},
	}
}

func (c *Cache) newID() string {
	id := strconv.Itoa(c.ID)
	c.ID++
	return id
}

func (c *Cache) Add(e Event) {
	e.ID = c.newID()
	c.Events[e.ID] = e
	year, month, day := e.Date.Date()

	if c.Tree[year] == nil {
		c.Tree[year] = make(map[time.Month]map[int][]string)
	}
	if c.Tree[year][month] == nil {
		c.Tree[year][month] = make(map[int][]string)
	}

	c.Tree[year][month][day] = append(c.Tree[year][month][day], e.ID)
}

func (c *Cache) Contains(e Event) bool {
	_, ok := c.Events[e.ID]
	return ok
}

func (c *Cache) Update(e Event) {
	dateOld := c.Events[e.ID].Date
	if dateOld.Equal(e.Date.Time) {
		c.Events[e.ID] = e
	} else {
		c.Delete(e)
		c.Add(e)
	}
}

func (c *Cache) Delete(e Event) {
	e = c.Events[e.ID]
	delete(c.Events, e.ID)
	year, month, day := e.Date.Date()
	events := c.Tree[year][month][day]
	for i, v := range events {
		if e.ID == v {
			events = append(events[:i], events[i+1:]...)
			break
		}
	}
	c.Tree[year][month][day] = events
}

func (c *Cache) getByID(id string) Event {
	return c.Events[id]
}

func (c *Cache) getByIDs(ids []string) (ev []Event) {
	for _, id := range ids {
		ev = append(ev, c.getByID(id))
	}
	return
}

func (c *Cache) GetByDate(date time.Time) (ev []Event) {
	year, month, day := date.Date()
	events := c.Tree[year][month][day]
	return c.getByIDs(events)
}

func (c *Cache) GetByWeek(date time.Time) (ev []Event) {
	for i := 0; i < 7; i++ {
		ev = append(ev, c.GetByDate(date)...)
		date = date.AddDate(0, 0, 1)
	}
	return
}

func (c *Cache) GetByMonth(date time.Time) (ev []Event) {
	year, month, _ := date.Date()
	monthEvents := c.Tree[year][month]
	for _, events := range monthEvents {
		ev = append(ev, c.getByIDs(events)...)
	}
	return
}
