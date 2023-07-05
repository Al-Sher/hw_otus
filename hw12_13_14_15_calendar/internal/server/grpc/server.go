package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	app    app.App
	config config.Config
	logger logger.Logger
	srv    *grpc.Server
	pb.CalendarServer
}

func NewServer(app app.App, l logger.Logger, c config.Config) *Server {
	return &Server{
		app:    app,
		config: c,
		logger: l,
	}
}

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.config.GRPCAddr())
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerRequestLoggerMiddlewareInterceptor(s.logger),
		),
	)

	pb.RegisterCalendarServer(s.srv, s)

	s.logger.Info(fmt.Sprintf("starting grpc server on %s", s.config.GRPCAddr()))

	err = s.srv.Serve(lsn)
	<-ctx.Done()

	return err
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}

func (s *Server) Create(ctx context.Context, e *pb.CreateEvent) (*pb.Result, error) {
	err := s.app.CreateEvent(
		ctx,
		e.GetTitle(),
		e.GetStartAt().AsTime(),
		e.GetDuration().AsDuration(),
		e.GetDescription(),
		e.GetAuthorId(),
	)
	if err != nil {
		return &pb.Result{}, status.Error(codes.Unknown, err.Error())
	}

	return &pb.Result{}, nil
}

func (s *Server) Update(ctx context.Context, e *pb.UpdateEvent) (*pb.Result, error) {
	err := s.app.UpdateEvent(
		ctx,
		e.GetId(),
		e.GetEvent().GetTitle(),
		e.GetEvent().GetStartAt().AsTime(),
		e.GetEvent().GetDuration().AsDuration(),
		e.GetEvent().GetDescription(),
		e.GetEvent().GetAuthorId(),
	)
	if err != nil {
		return &pb.Result{}, status.Error(codes.Unknown, err.Error())
	}

	return &pb.Result{}, nil
}

func (s *Server) Delete(ctx context.Context, e *pb.DeleteEvent) (*pb.Result, error) {
	err := s.app.DeleteEvent(
		ctx,
		e.GetId(),
	)
	if err != nil {
		return &pb.Result{}, status.Error(codes.Unknown, err.Error())
	}

	return &pb.Result{}, nil
}

func (s *Server) EventByDay(ctx context.Context, e *pb.EventDay) (*pb.EventsResult, error) {
	events, err := s.app.EventByDay(
		ctx,
		e.GetDate().AsTime(),
	)
	if err != nil {
		return &pb.EventsResult{}, status.Error(codes.Unknown, err.Error())
	}

	return convert(events), nil
}

func (s *Server) EventByWeek(ctx context.Context, e *pb.EventDay) (*pb.EventsResult, error) {
	events, err := s.app.EventByWeek(
		ctx,
		e.GetDate().AsTime(),
	)
	if err != nil {
		return &pb.EventsResult{}, status.Error(codes.Unknown, err.Error())
	}

	return convert(events), nil
}

func (s *Server) EventByMonth(ctx context.Context, e *pb.EventDay) (*pb.EventsResult, error) {
	events, err := s.app.EventByMonth(
		ctx,
		e.GetDate().AsTime(),
	)
	if err != nil {
		return &pb.EventsResult{}, status.Error(codes.Unknown, err.Error())
	}

	return convert(events), nil
}

func convert(events []storage.Event) *pb.EventsResult {
	result := make([]*pb.Event, 0, len(events))
	for _, event := range events {
		eventResult := &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			StartAt:     timestamppb.New(event.StartAt),
			Duration:    durationpb.New(event.EndAt.Sub(event.StartAt)),
			Description: event.Description,
			AuthorId:    event.AuthorID,
		}
		result = append(result, eventResult)
	}

	return &pb.EventsResult{Events: result}
}
