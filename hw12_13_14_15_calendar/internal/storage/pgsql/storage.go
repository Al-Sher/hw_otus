package sqlstorage

import (
	"context"
	"errors"
	"time"

	internalStorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5"
)

const Type string = "pgsql"

type storage struct {
	conn *pgx.Conn
}

func New() internalStorage.Storage {
	return &storage{}
}

func (s *storage) Connect(ctx context.Context, dsn string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	s.conn = conn

	return nil
}

func (s *storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *storage) CreateEvent(ctx context.Context, event internalStorage.Event) error {
	isBusy, err := s.isDateBusy(ctx, event)
	if err != nil {
		return err
	}

	if isBusy {
		return internalStorage.ErrDateBusy
	}

	sql := `INSERT INTO events 
    (id, title, start_at, end_at, description, author_id, notification_date) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = s.conn.Exec(
		ctx,
		sql,
		event.ID,
		event.Title,
		event.StartAt,
		event.EndAt,
		event.Description,
		event.AuthorID,
		event.NotificationDate,
	)

	return err
}

func (s *storage) UpdateEvent(ctx context.Context, event internalStorage.Event) error {
	isExist, err := s.isExistByID(ctx, event.ID)
	if err != nil {
		return err
	}
	if !isExist {
		return internalStorage.ErrEventNotFound
	}

	isBusy, err := s.isDateBusy(ctx, event)
	if err != nil {
		return err
	}

	if isBusy {
		return internalStorage.ErrDateBusy
	}

	sql := `UPDATE events 
	SET title=$2, start_at=$3, end_at=$4, description=$5, author_id=$6, notification_date=$7 
	WHERE id = $1`

	_, err = s.conn.Exec(
		ctx,
		sql,
		event.ID,
		event.Title,
		event.StartAt,
		event.EndAt,
		event.Description,
		event.AuthorID,
		event.NotificationDate,
	)

	return err
}

func (s *storage) DeleteEvent(ctx context.Context, id string) error {
	isExist, err := s.isExistByID(ctx, id)
	if err != nil {
		return err
	}
	if !isExist {
		return internalStorage.ErrEventNotFound
	}

	sql := `DELETE FROM events WHERE id=$1`

	_, err = s.conn.Exec(ctx, sql, id)

	return err
}

func (s *storage) EventsDay(ctx context.Context, date time.Time) ([]internalStorage.Event, error) {
	year, month, day := date.Date()
	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, month, day, 23, 59, 59, 999, time.UTC)

	return s.eventsByDates(ctx, startDate, endDate)
}

func (s *storage) EventsWeek(ctx context.Context, date time.Time) ([]internalStorage.Event, error) {
	year, month, day := date.Date()
	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, month, day+7, 23, 59, 59, 999, time.UTC)

	return s.eventsByDates(ctx, startDate, endDate)
}

func (s *storage) EventsMonth(ctx context.Context, date time.Time) ([]internalStorage.Event, error) {
	year, month, day := date.Date()
	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, month+1, day, 23, 59, 59, 999, time.UTC)

	return s.eventsByDates(ctx, startDate, endDate)
}

func (s *storage) EventsForNotification(ctx context.Context) ([]internalStorage.Event, error) {
	result := make([]internalStorage.Event, 0)

	sql := `SELECT id, title, start_at, end_at, description, author_id, notification_date 
	FROM events 
	WHERE notification_date IS NOT NULL AND notification_date < NOW()`

	rows, err := s.conn.Query(ctx, sql)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		event := internalStorage.Event{}
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartAt,
			&event.EndAt,
			&event.Description,
			&event.AuthorID,
			&event.NotificationDate,
		); err != nil {
			return nil, err
		}

		result = append(result, event)
	}

	return result, nil
}

func (s *storage) ClearNotificationDates(ctx context.Context, ids []string) error {
	sql := `UPDATE events SET notification_date = NULL WHERE id = ANY($1)`

	_, err := s.conn.Exec(ctx, sql, ids)
	return err
}

func (s *storage) ClearOldEvents(ctx context.Context) error {
	sql := `DELETE FROM events WHERE start_at < NOW()- interval '1 year'`

	_, err := s.conn.Exec(ctx, sql)
	return err
}

func (s *storage) eventsByDates(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]internalStorage.Event, error) {
	result := make([]internalStorage.Event, 0)

	sql := `SELECT id, title, start_at, end_at, description, author_id, notification_date 
	FROM events 
	WHERE start_at between $1 and $2`
	rows, err := s.conn.Query(ctx, sql, startDate, endDate)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		event := internalStorage.Event{}
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartAt,
			&event.EndAt,
			&event.Description,
			&event.AuthorID,
			&event.NotificationDate,
		); err != nil {
			return nil, err
		}

		result = append(result, event)
	}

	return result, nil
}

func (s *storage) isDateBusy(ctx context.Context, event internalStorage.Event) (bool, error) {
	sql := `SELECT id from events WHERE (start_at BETWEEN $1 AND $2 OR end_at BETWEEN $1 AND $2) AND id != $3`
	row := s.conn.QueryRow(ctx, sql, event.StartAt, event.EndAt, event.ID)

	var id string
	err := row.Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	return true, err
}

func (s *storage) isExistByID(ctx context.Context, id string) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`
	row := s.conn.QueryRow(ctx, sql, id)
	isExist := false

	err := row.Scan(&isExist)
	return isExist, err
}
