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
		NotificationAt time.Time,
	) error
	UpdateEvent(
		ctx context.Context,
		id string,
		title string,
		startAt time.Time,
		duration time.Duration,
		description string,
		authorID string,
		NotificationAt time.Time,
	) error
	DeleteEvent(ctx context.Context, id string) error
	EventByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	EventByWeek(ctx context.Context, day time.Time) ([]storage.Event, error)
	EventByMonth(ctx context.Context, day time.Time) ([]storage.Event, error)
	EventsForNotification(ctx context.Context) ([]storage.Notification, error)
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
	notificationAt time.Time,
) error {
	id := uuid.NewString()
	return a.storage.CreateEvent(ctx, storage.Event{
		ID:               id,
		Title:            title,
		StartAt:          startAt,
		EndAt:            startAt.Add(duration),
		Description:      description,
		AuthorID:         authorID,
		NotificationDate: notificationAt,
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
	notificationAt time.Time,
) error {
	return a.storage.UpdateEvent(ctx, storage.Event{
		ID:               id,
		Title:            title,
		StartAt:          startAt,
		EndAt:            startAt.Add(duration),
		Description:      description,
		AuthorID:         authorID,
		NotificationDate: notificationAt,
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

func (a *app) EventsForNotification(ctx context.Context) ([]storage.Notification, error) {
	events, err := a.storage.EventsForNotification(ctx)
	if err != nil {
		return nil, err
	}
	notifications := make([]storage.Notification, len(events), 0)

	for _, event := range events {
		notification := storage.Notification{
			ID:       event.ID,
			Title:    event.Title,
			Date:     event.NotificationDate,
			AuthorID: event.AuthorID,
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
