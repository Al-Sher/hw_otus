package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

const exceptedEmpty = "{}"

func TestHandler_Handlers(t *testing.T) {
	s := memorystorage.New()
	a := app.New(s)
	l, err := logger.New("debug", "/dev/stdout")
	require.NoError(t, err)
	ctx := context.Background()

	d, err := time.Parse(time.DateOnly, "2023-06-01")
	require.NoError(t, err)
	e := event{
		Title:       "Test event",
		StartAt:     d,
		Duration:    3600,
		Description: "Test description",
		AuthorID:    "512b922c-822a-4a05-b52b-85b85ab7a00c",
	}

	h := NewHandlers(a, l)

	httpClient := &http.Client{}

	t.Run("Empty Request", func(t *testing.T) {
		test := httptest.NewServer(h.Handlers(ctx))
		defer test.Close()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, test.URL+"/events/day/2000-01-01", nil)
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		out, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, exceptedEmpty, string(out))
	})

	t.Run("Not found", func(t *testing.T) {
		test := httptest.NewServer(h.Handlers(ctx))
		defer test.Close()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, test.URL+"/events/2000-01-01", nil)
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("basic", func(t *testing.T) {
		test := httptest.NewServer(h.Handlers(ctx))
		defer test.Close()
		data, err := json.Marshal(e)
		require.NoError(t, err)

		// Создадим новую запись
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, test.URL+"/events", bytes.NewReader(data))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		// Проверим что запись создалась
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, test.URL+"/events/day/"+d.Format(time.DateOnly), nil)
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		out, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		// Проверим что запись имеет правильные значения и сохраним сгенерированный id
		te := &result{}
		err = json.Unmarshal(out, te)
		require.NoError(t, err)
		require.Equal(t, 1, len(te.Events))
		require.Equal(t, e.Title, te.Events[0].Title)
		require.Equal(t, e.StartAt, te.Events[0].StartAt)
		require.Equal(t, e.Duration, te.Events[0].Duration)
		require.Equal(t, e.Description, te.Events[0].Description)
		require.Equal(t, e.AuthorID, te.Events[0].AuthorID)
		id := te.Events[0].ID
		// Обновим запись
		e.Title = "Test after update"
		data, err = json.Marshal(e)
		require.NoError(t, err)
		req, err = http.NewRequestWithContext(ctx, http.MethodPut, test.URL+"/events/"+id, bytes.NewReader(data))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		// Проверим что запись обновилась
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, test.URL+"/events/day/"+d.Format(time.DateOnly), nil)
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		out, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		// Проверим что запись имеет правильные значения
		err = json.Unmarshal(out, &te)
		require.NoError(t, err)
		require.Equal(t, e.Title, te.Events[0].Title)
		require.Equal(t, e.StartAt, te.Events[0].StartAt)
		require.Equal(t, e.Duration, te.Events[0].Duration)
		require.Equal(t, e.Description, te.Events[0].Description)
		require.Equal(t, e.AuthorID, te.Events[0].AuthorID)
		require.Equal(t, id, te.Events[0].ID)
		// Удалим запись
		req, err = http.NewRequestWithContext(ctx, http.MethodDelete, test.URL+"/events/"+id, nil)
		require.NoError(t, err)
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		// Проверим что запись удалилась
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, test.URL+"/events/day/"+d.Format(time.DateOnly), nil)
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/json")
		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		out, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, exceptedEmpty, string(out))
	})
}
