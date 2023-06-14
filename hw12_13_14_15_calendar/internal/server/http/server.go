package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	app app.App
	srv *http.Server
}

func NewServer(app app.App) *Server {
	return &Server{
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.app.Logger().Info("Start server on " + s.app.Config().HTTPAddr())
	s.srv = &http.Server{
		Addr:              s.app.Config().HTTPAddr(),
		Handler:           handler(ctx, s.app),
		ReadHeaderTimeout: time.Duration(s.app.Config().HTTPReadTimeout()) * time.Second,
	}

	err := s.srv.ListenAndServe()
	<-ctx.Done()

	// Завершение работы сервера не является ошибкой запуска
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func handler(ctx context.Context, app app.App) http.Handler {
	service := NewHandlers(app)

	h := loggingMiddleware(ctx, service.Handlers(ctx), app.Logger())

	return h
}
