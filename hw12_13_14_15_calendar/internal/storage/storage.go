package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDateBusy      = errors.New("данное время уже занято другим событием")
	ErrEventNotFound = errors.New("событие не найдено")
)

type Storage interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id string) error
	EventsDay(ctx context.Context, date time.Time) ([]Event, error)
	EventsWeek(ctx context.Context, date time.Time) ([]Event, error)
	EventsMonth(ctx context.Context, date time.Time) ([]Event, error)
	EventsByDates(ctx context.Context, startDate time.Time, endDate time.Time) ([]Event, error)
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
}