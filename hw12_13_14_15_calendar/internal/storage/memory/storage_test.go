package memorystorage

import (
	"context"
	"errors"
	"testing"
	"time"

	internalStorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	startDateTime1, err := time.Parse(time.RFC3339, "2023-06-01T21:00:00+03:00")
	require.NoError(t, err)
	startDateTime2, err := time.Parse(time.RFC3339, "2023-06-02T21:00:00+03:00")
	require.NoError(t, err)
	startDateTime3, err := time.Parse(time.RFC3339, "2023-06-01T21:30:00+03:00")
	require.NoError(t, err)
	startDateTime4, err := time.Parse(time.RFC3339, "2023-06-03T21:30:00+03:00")
	require.NoError(t, err)
	events := []internalStorage.Event{
		{
			ID:          "1",
			Title:       "test",
			Description: "test description",
			StartAt:     startDateTime1,
			EndAt:       startDateTime1.Add(1 * time.Hour),
			AuthorID:    "1",
		},
		{
			ID:          "2",
			Title:       "test",
			Description: "test description",
			StartAt:     startDateTime2,
			EndAt:       startDateTime2.Add(1 * time.Hour),
			AuthorID:    "1",
		},
	}
	errEvent := internalStorage.Event{
		ID:          "3",
		Title:       "test",
		Description: "test description",
		StartAt:     startDateTime3,
		EndAt:       startDateTime3.Add(10 * time.Minute),
		AuthorID:    "1",
	}
	s := New()
	t.Run("create events", func(t *testing.T) {
		for _, v := range events {
			err := s.CreateEvent(ctx, v)
			require.NoError(t, err)
		}
	})
	t.Run("create event with error", func(t *testing.T) {
		err := s.CreateEvent(ctx, errEvent)
		require.Truef(
			t,
			errors.Is(err, internalStorage.ErrDateBusy),
			"actual error %q, excepted %q",
			err,
			internalStorage.ErrDateBusy,
		)
	})
	t.Run("get events by day", func(t *testing.T) {
		events, err := s.EventsDay(ctx, startDateTime1)
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
	})
	t.Run("get events by week", func(t *testing.T) {
		events, err := s.EventsWeek(ctx, startDateTime1)
		require.NoError(t, err)
		require.Equal(t, 2, len(events))
	})
	t.Run("update event", func(t *testing.T) {
		event := internalStorage.Event{
			ID:          "1",
			Title:       "test with update",
			Description: "test description",
			StartAt:     startDateTime1,
			EndAt:       startDateTime1.Add(1 * time.Hour),
			AuthorID:    "1",
		}
		err := s.UpdateEvent(ctx, event)
		require.NoError(t, err)
		events, err := s.EventsDay(ctx, startDateTime1)
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0], event)
	})
	t.Run("update event date time without error", func(t *testing.T) {
		event := internalStorage.Event{
			ID:          "1",
			Title:       "test with update",
			Description: "test description",
			StartAt:     startDateTime4,
			EndAt:       startDateTime4.Add(1 * time.Hour),
			AuthorID:    "1",
		}
		err := s.UpdateEvent(ctx, event)
		require.NoError(t, err)

		events, err := s.EventsDay(ctx, startDateTime4)
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0], event)
	})
	t.Run("update not found event", func(t *testing.T) {
		event := internalStorage.Event{
			ID:          "5",
			Title:       "test with update",
			Description: "test description",
			StartAt:     startDateTime4,
			EndAt:       startDateTime4.Add(1 * time.Hour),
			AuthorID:    "1",
		}
		err := s.UpdateEvent(ctx, event)
		require.Truef(
			t,
			errors.Is(err, internalStorage.ErrEventNotFound),
			"actual error %q, excepted %q",
			err,
			internalStorage.ErrDateBusy,
		)
	})
	t.Run("update with ErrDateBusy", func(t *testing.T) {
		event := internalStorage.Event{
			ID:          "1",
			Title:       "test with update",
			Description: "test description",
			StartAt:     startDateTime2,
			EndAt:       startDateTime2.Add(1 * time.Hour),
			AuthorID:    "1",
		}
		err := s.UpdateEvent(ctx, event)
		require.Truef(
			t,
			errors.Is(err, internalStorage.ErrDateBusy),
			"actual error %q, excepted %q",
			err,
			internalStorage.ErrDateBusy,
		)
	})
	t.Run("delete event", func(t *testing.T) {
		idForDelete := "2"
		err := s.DeleteEvent(ctx, idForDelete)
		require.NoError(t, err)
		events, err := s.EventsWeek(ctx, startDateTime1)
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.NotEqual(t, events[0].ID, idForDelete)
	})
	t.Run("delete not found event", func(t *testing.T) {
		err := s.DeleteEvent(ctx, "2")
		require.Truef(
			t,
			errors.Is(err, internalStorage.ErrEventNotFound),
			"actual error %q, excepted %q",
			err,
			internalStorage.ErrDateBusy,
		)
	})
}
