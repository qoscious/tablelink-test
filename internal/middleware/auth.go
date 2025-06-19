package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TODO: Implementasi middleware auth

func AuthInterceptor(jwtSecret string, section string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, grpc.Errorf(16, "missing metadata")
		}
		// Cek header X-Link-Service
		if len(md["x-link-service"]) == 0 || md["x-link-service"][0] != section {
			return nil, grpc.Errorf(16, "invalid section header")
		}
		// Cek Authorization Bearer
		tokens := md["authorization"]
		if len(tokens) == 0 || !strings.HasPrefix(tokens[0], "Bearer ") {
			return nil, grpc.Errorf(16, "missing or invalid token")
		}
		// TODO: Validasi JWT dan cek di Redis
		return handler(ctx, req)
	}
}
