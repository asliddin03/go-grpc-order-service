package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	st, _ := status.FromError(err)
	log.Printf(
		"grpc method=%s code=%s duration=%s req=%v",
		info.FullMethod,
		st.Code(),
		time.Since(start),
		err,
	)

	return resp, err
}
