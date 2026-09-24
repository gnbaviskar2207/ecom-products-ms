package grpc

import (
	"context"
	"errors"
	"log/slog"

	errs "github.com/gnbaviskar2207/ecom-common/pkg/err"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled by client")
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		}
		return nil, errs.ToGRPCStatusError(ctx, err, logger, info.FullMethod)
	}
}
