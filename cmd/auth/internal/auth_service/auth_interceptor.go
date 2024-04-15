package authservice

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTManager struct {
	secretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey}
}

type AuthInterceptor struct {
	jwtManager *JWTManager
}

func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (jwtManager *JWTManager) VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing
		return []byte(jwtManager.secretKey), nil
	})
	println(token.Valid, "this is token value")

	if err != nil {
		return nil, err
	}
	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	// Extract and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	// TODO: implement this
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md, "this is md")
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	_, err := interceptor.jwtManager.VerifyJWT(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	return nil
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)
		if info.FullMethod == "/auth.AuthService/LoginUser" || info.FullMethod == "/auth.AuthService/RegisterUser" {
			return handler(ctx, req)
		}

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
