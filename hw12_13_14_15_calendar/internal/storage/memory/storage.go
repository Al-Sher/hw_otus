package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	internalStorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
)

const (
	Type       string = "inMemory"
	dateLayout string = "2006-01-02"
)

type storage struct {
	events         map[string]internalStorage.Event
	eventIdsByDate map[string]map[string]struct{}
	sendingEvents  map[string]struct{}
	mu             sync.RWMutex
}

func New() internalStorage.Storage {
	return &storage{
		mu:             sync.RWMutex{},
		events:         make(map[string]internalStorage.Event),
		eventIdsByDate: make(map[string]map[string]struct{}),
		sendingEvents:  make(map[string]struct{}),
	}
}

func (s *storage) CreateEvent(_ context.Context, event internalStorage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	year, month, day := event.StartAt.Date()
	dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)

	if _, ok := s.eventIdsByDate[dateStr]; !ok {
		s.eventIdsByDate[dateStr] = make(map[string]struct{})
	}

	if s.isDateBusy(event) {
		return internalStorage.ErrDateBusy
	}

	s.events[event.ID] = event
	s.eventIdsByDate[dateStr][event.ID] = struct{}{}

	return nil
}

func (s *storage) UpdateEvent(_ context.Context, event internalStorage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldEvent, ok := s.events[event.ID]
	if !ok {
		return internalStorage.ErrEventNotFound
	}

	if s.isDateBusy(event) {
		return internalStorage.ErrDateBusy
	}

	if oldEvent.StartAt != event.StartAt {
		year, month, day := oldEvent.StartAt.Date()
		dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)

		delete(s.eventIdsByDate[dateStr], oldEvent.ID)

		year, month, day = event.StartAt.Date()
		dateStr = fmt.Sprintf("%d-%d-%d", year, month, day)

		if _, ok := s.eventIdsByDate[dateStr]; !ok {
			s.eventIdsByDate[dateStr] = make(map[string]struct{})
		}

		s.eventIdsByDate[dateStr][event.ID] = struct{}{}
	}

	s.events[event.ID] = event

	return nil
}

func (s *storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	eventForDelete, ok := s.events[id]
	if !ok {
		return internalStorage.ErrEventNotFound
	}

	year, month, day := eventForDelete.StartAt.Date()
	dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)
	delete(s.eventIdsByDate[dateStr], id)
	delete(s.events, id)

	return nil
}

func (s *storage) EventsDay(ctx context.Context, date time.Time) ([]internalStorage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.eventsByDates(ctx, date, date.AddDate(0, 0, 1))
}

func (s *storage) EventsWeek(ctx context.Context, date time.Time) ([]internalStorage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.eventsByDates(ctx, date, date.AddDate(0, 0, 7))
}

func (s *storage) EventsMonth(ctx context.Context, date time.Time) ([]internalStorage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.eventsByDates(ctx, date, date.AddDate(0, 1, 0))
}

func (s *storage) eventsByDates(
	_ context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]internalStorage.Event, error) {
	days := endDate.Sub(startDate).Hours() / 24

	result := make([]internalStorage.Event, 0)
	for i := 0; i < int(days); i++ {
		t := startDate.AddDate(0, 0, i)

		year, month, day := t.Date()
		dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)

		for k := range s.eventIdsByDate[dateStr] {
			result = append(result, s.events[k])
		}
	}

	return result, nil
}

func (s *storage) EventsForNotification(_ context.Context) ([]internalStorage.Event, error) {
	res := make([]internalStorage.Event, 0)
	for _, event := range s.events {
		if _, ok := s.sendingEvents[event.ID]; !ok && event.NotificationDate.Before(time.Now()) {
			res = append(res, event)
		}
	}

	return res, nil
}

func (s *storage) ClearNotificationDates(_ context.Context, ids []string) error {
	for _, id := range ids {
		s.sendingEvents[id] = struct{}{}
	}

	return nil
}

func (s *storage) ClearOldEvents(_ context.Context) error {
	dateClear := time.Now().AddDate(-1, 0, 0)
	for dateStr, events := range s.eventIdsByDate {
		t, err := time.Parse(dateLayout, dateStr)
		if err != nil {
			return err
		}
		if t.Before(dateClear) {
			for id := range events {
				delete(s.events, id)
				delete(s.sendingEvents, id)
			}
			delete(s.eventIdsByDate, dateStr)
		}
	}

	return nil
}

func (s *storage) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *storage) Close(_ context.Context) error {
	return nil
}

func (s *storage) isDateBusy(event internalStorage.Event) bool {
	year, month, day := event.StartAt.Date()
	dateStr := fmt.Sprintf("%d-%d-%d", year, month, day)

	eventStart := event.StartAt.Unix()
	eventEnd := event.EndAt.Unix()

	for id := range s.eventIdsByDate[dateStr] {
		t := s.events[id]
		if t.EndAt.Unix() <= eventStart || t.StartAt.Unix() >= eventEnd || t.ID == event.ID {
			continue
		}

		return true
	}

	return false
}
