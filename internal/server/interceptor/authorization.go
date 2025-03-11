package interceptor

import (
	"context"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/pb"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/server/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var whitelistedMethods = []string{
	pb.UserService_Login_FullMethodName,
	pb.UserService_Register_FullMethodName,
}

const (
	JWTHeader        = "jwt"
	UserIDContextKey = "UserID"
)

func AuthorizationInterceptor(jwtService *security.JwtService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip Register and Login methods
		for _, method := range whitelistedMethods {
			if strings.EqualFold(info.FullMethod, method) {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		jwtMetadata, exists := md[JWTHeader]
		if !exists || len(jwtMetadata) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing or malformed JWT")
		}

		jwt := jwtMetadata[0]

		claims, err := jwtService.ValidateJwtToken(jwt)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)

		return handler(ctx, req)
	}
}
