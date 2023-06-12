package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
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
		Handler:           mux(ctx, s.app.Logger()),
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

func mux(ctx context.Context, logger logger.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.RequestURI != "/" {
			if err := notFound(writer); err != nil {
				logger.Error(err)
				return
			}
			return
		}

		_, err := writer.Write([]byte("hello-world"))
		if err != nil {
			logger.Error(err)
			return
		}
	})

	handler := loggingMiddleware(ctx, mux, logger)

	return handler
}

func notFound(writer http.ResponseWriter) error {
	writer.WriteHeader(http.StatusNotFound)
	_, err := writer.Write([]byte("Page not found"))
	return err
}
