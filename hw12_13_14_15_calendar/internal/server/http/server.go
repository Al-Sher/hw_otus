package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	app    app.App
	srv    *http.Server
	logger logger.Logger
	config config.Config
}

func NewServer(app app.App, l logger.Logger, c config.Config) *Server {
	return &Server{
		app:    app,
		logger: l,
		config: c,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Start server on " + s.config.HTTPAddr())
	s.srv = &http.Server{
		Addr:              s.config.HTTPAddr(),
		Handler:           mux(ctx, s.logger),
		ReadHeaderTimeout: time.Duration(s.config.HTTPReadTimeout()) * time.Second,
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
