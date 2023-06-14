package grpc

import (
	"fmt"
	"time"

	"github.com/Al-Sher/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

import (
	"context"
)

func UnaryServerRequestLoggerMiddlewareInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, r interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		result, err := handler(ctx, r)

		var realIp string
		if p, ok := peer.FromContext(ctx); ok {
			realIp = p.Addr.String()
		}

		var userAgent string
		if p, ok := metadata.FromIncomingContext(ctx); ok {
			userAgents := p.Get("user-agent")
			if len(userAgents) > 0 {
				userAgent = userAgents[0]
			}
		}

		logger.Info(
			fmt.Sprintf(
				"%s [%s] %s %s %s %d %f \"%s\"",
				realIp,
				time.Now(),
				"grpc",
				info.FullMethod,
				"grpc",
				status.Code(err),
				time.Since(start).Seconds(),
				userAgent,
			),
		)

		return result, err
	}
}
