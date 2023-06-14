package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
)

type Handler struct {
	app    app.App
	logger logger.Logger
}

type result struct {
	Events  []*event `json:"events,omitempty"`
	Error   error    `json:"error,omitempty"`
	Success string   `json:"success,omitempty"`
}

type event struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title"`
	StartAt     time.Time `json:"startAt"`
	Duration    float64   `json:"duration"`
	Description string    `json:"description"`
	AuthorID    string    `json:"authorId"`
}

type EventResult struct {
	Events []*event `json:"events"`
}

const URLPath = "/events"

var (
	ErrNotSupportedMethod = errors.New("unsupported method")
	ErrPageNotFound       = errors.New("page not found")
)

func NewHandlers(app app.App, l logger.Logger) Handler {
	return Handler{app: app, logger: l}
}

func (h *Handler) Handlers(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, URLPath) {
			if err := notFound(w); err != nil {
				h.logger.Error(err)
				return
			}
		}

		var res result
		switch r.Method {
		case http.MethodPost:
			res = h.create(ctx, r)
		case http.MethodPut:
			res = h.update(ctx, r)
		case http.MethodDelete:
			res = h.delete(ctx, r)
		case http.MethodGet:
			res = h.get(ctx, r)
		default:
			res = result{Error: ErrNotSupportedMethod}
		}

		data, err := json.Marshal(res)
		if err != nil {
			h.logger.Error(err)
			if err := internalError(w, err); err != nil {
				h.logger.Error(err)
			}
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if res.Error != nil {
			switch {
			case errors.Is(res.Error, ErrPageNotFound):
				w.WriteHeader(http.StatusNotFound)
			case errors.Is(res.Error, ErrNotSupportedMethod):
				w.WriteHeader(http.StatusMethodNotAllowed)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}

			h.logger.Error(res.Error)
		}

		if _, err := w.Write(data); err != nil {
			h.logger.Error(err)
		}
	}
}

func notFound(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Page not found"))
	return err
}

func internalError(w http.ResponseWriter, err error) error {
	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write([]byte(err.Error()))
	return e
}

func (h *Handler) unmarshalEvent(r *http.Request) (*event, error) {
	content, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	e := &event{}
	if err := json.Unmarshal(content, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (h *Handler) create(ctx context.Context, r *http.Request) result {
	e, err := h.unmarshalEvent(r)
	if err != nil {
		return result{Error: err}
	}

	if err := h.app.CreateEvent(
		ctx,
		e.Title,
		e.StartAt,
		time.Duration(e.Duration)*time.Second,
		e.Description,
		e.AuthorID,
	); err != nil {
		return result{Error: err}
	}

	return result{Success: "Событие успешно добавлено в календарь"}
}

func (h *Handler) update(ctx context.Context, r *http.Request) result {
	e, err := h.unmarshalEvent(r)
	if err != nil {
		return result{Error: err}
	}

	if r.URL.Path == URLPath {
		return result{Error: ErrNotSupportedMethod}
	}

	id := r.URL.Path[len(URLPath)+1:]

	if err := h.app.UpdateEvent(
		ctx,
		id,
		e.Title,
		e.StartAt,
		time.Duration(e.Duration)*time.Second,
		e.Description,
		e.AuthorID,
	); err != nil {
		return result{Error: err}
	}

	return result{Success: "Событие успешно обновлено"}
}

func (h *Handler) delete(ctx context.Context, r *http.Request) result {
	if r.URL.Path == URLPath {
		return result{Error: ErrNotSupportedMethod}
	}

	id := r.URL.Path[len(URLPath)+1:]

	if err := h.app.DeleteEvent(
		ctx,
		id,
	); err != nil {
		return result{Error: err}
	}

	return result{Success: "Событие успешно удалено"}
}

func (h *Handler) get(ctx context.Context, r *http.Request) result {
	if r.URL.Path == URLPath {
		return result{Error: ErrNotSupportedMethod}
	}
	query := strings.Split(r.URL.Path[len(URLPath)+1:], "/")
	if len(query) != 2 || !contains(query[0], []string{"day", "week", "month"}) {
		return result{Error: ErrPageNotFound}
	}

	method := query[0]
	t, err := time.Parse(time.DateOnly, query[1])
	if err != nil {
		return result{Error: err}
	}

	switch method {
	case "day":
		return h.eventByDay(ctx, t)
	case "week":
		return h.eventByWeek(ctx, t)
	case "month":
		return h.eventByMonth(ctx, t)
	default:
		return result{Error: ErrPageNotFound}
	}
}

func (h *Handler) eventByDay(ctx context.Context, d time.Time) result {
	events, err := h.app.EventByDay(
		ctx,
		d,
	)

	return result{Error: err, Events: convert(events)}
}

func (h *Handler) eventByWeek(ctx context.Context, d time.Time) result {
	events, err := h.app.EventByWeek(
		ctx,
		d,
	)

	return result{Error: err, Events: convert(events)}
}

func (h *Handler) eventByMonth(ctx context.Context, d time.Time) result {
	events, err := h.app.EventByMonth(
		ctx,
		d,
	)

	return result{Error: err, Events: convert(events)}
}

func convert(events []storage.Event) []*event {
	r := make([]*event, 0, len(events))
	for _, e := range events {
		eResult := &event{
			ID:          e.ID,
			Title:       e.Title,
			StartAt:     e.StartAt,
			Duration:    e.EndAt.Sub(e.StartAt).Seconds(),
			Description: e.Description,
			AuthorID:    e.AuthorID,
		}
		r = append(r, eResult)
	}

	return r
}

func contains(s string, a []string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}

	return false
}
