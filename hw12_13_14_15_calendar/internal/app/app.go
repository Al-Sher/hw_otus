package app

import (
	"context"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App interface {
	CreateEvent(
		ctx context.Context,
		title string,
		startAt time.Time,
		duration time.Duration,
		description string,
		authorID string,
	) error
	UpdateEvent(
		ctx context.Context,
		id string,
		title string,
		startAt time.Time,
		duration time.Duration,
		description string,
		authorID string,
	) error
	DeleteEvent(ctx context.Context, id string) error
	EventByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	EventByWeek(ctx context.Context, day time.Time) ([]storage.Event, error)
	EventByMonth(ctx context.Context, day time.Time) ([]storage.Event, error)
}

type app struct {
	storage storage.Storage
}

func New(storage storage.Storage) App {
	return &app{
		storage: storage,
	}
}

func (a *app) CreateEvent(
	ctx context.Context,
	title string,
	startAt time.Time,
	duration time.Duration,
	description string,
	authorID string,
) error {
	id := uuid.NewString()

	return a.storage.CreateEvent(ctx, storage.Event{
		ID:          id,
		Title:       title,
		StartAt:     startAt,
		EndAt:       startAt.Add(duration),
		Description: description,
		AuthorID:    authorID,
	})
}

func (a *app) UpdateEvent(
	ctx context.Context,
	id string,
	title string,
	startAt time.Time,
	duration time.Duration,
	description string,
	authorID string,
) error {
	return a.storage.UpdateEvent(ctx, storage.Event{
		ID:          id,
		Title:       title,
		StartAt:     startAt,
		EndAt:       startAt.Add(duration),
		Description: description,
		AuthorID:    authorID,
	})
}

func (a *app) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *app) EventByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	return a.storage.EventsDay(ctx, day)
}

func (a *app) EventByWeek(ctx context.Context, day time.Time) ([]storage.Event, error) {
	return a.storage.EventsWeek(ctx, day)
}

func (a *app) EventByMonth(ctx context.Context, day time.Time) ([]storage.Event, error) {
	return a.storage.EventsMonth(ctx, day)
}
