package interceptor

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRecoveryInterceptor(ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf(
				"panic recovered method=%s panic=%v\n%s",
				info.FullMethod,
				r,
				debug.Stack(),
			)

			err = status.Error(codes.Internal, "internal server error")
		}
	}()

	return handler(ctx, req)
}
