-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
                        id uuid primary key,
                        title text,
                        start_at timestamp,
                        end_at timestamp,
                        description text,
                        author_id uuid,
                        notification_date timestamp
);

CREATE INDEX ix_events_date ON events (start_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
