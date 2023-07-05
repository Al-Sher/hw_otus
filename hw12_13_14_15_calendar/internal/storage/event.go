package storage

import "time"

type Event struct {
	ID               string
	Title            string
	StartAt          time.Time
	EndAt            time.Time
	Description      string
	AuthorID         string
	NotificationDate time.Time
}

type Notification struct {
	ID       string
	Title    string
	Date     time.Time
	AuthorID string
}
