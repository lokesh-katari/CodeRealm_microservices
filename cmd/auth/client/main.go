package main

import (
	"context"
	"log"
	"time"

	"lokesh-katari/code-realm/cmd/auth/internal/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthClient struct {
	service      auth.AuthServiceClient
	username     string
	password     string
	initialToken string
}

type AuthInterceptor struct {
	authClient   *AuthClient
	accessToken  string
	initialToken string
}

func NewAuthInterceptor(authClient *AuthClient, refreshDuration time.Duration, initialToken string) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:   authClient,
		accessToken:  initialToken,
		initialToken: initialToken,
	}

	// err := interceptor.refreshToken(refreshDuration)
	// if err != nil {
	//     return nil, err
	// }

	return interceptor, nil
}

func (interceptor *AuthInterceptor) refreshToken() error {
	accessToken, err := interceptor.authClient.Login()
	if err != nil {
		return err
	}
	interceptor.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)
	return nil
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	if interceptor.accessToken == "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.initialToken)
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

func NewAuthClient(cc *grpc.ClientConn, username, password, initialToken string) *AuthClient {
	service := auth.NewAuthServiceClient(cc)
	return &AuthClient{
		service,
		username,
		password,
		initialToken,
	}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &auth.LoginUserRequest{
		Email:    client.username,
		Password: client.password,
	}

	res, err := client.service.LoginUser(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetToken(), nil
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("--> unary interceptor: %s", method)
		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

func InitialLogin(cc *grpc.ClientConn, username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := auth.NewAuthServiceClient(cc)
	req := &auth.LoginUserRequest{
		Email:    username,
		Password: password,
	}

	res, err := client.LoginUser(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetToken(), nil
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	initialToken, err := InitialLogin(conn, "lokeshkatari921@outlook.com", "asdfasdf")
	if err != nil {
		panic(err)
	}

	authClient := NewAuthClient(conn, "lokeshkatari921@outlook.com", "asdfasdf", initialToken)
	interceptor, err := NewAuthInterceptor(authClient, 5*time.Second, initialToken)
	if err != nil {
		panic(err)
	}

	cc2, err := grpc.Dial(
		"localhost:50051",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
	)
	if err != nil {
		panic(err)
	}
	defer cc2.Close()

	client2 := auth.NewAuthServiceClient(cc2)
	res2, err := client2.UpdateUserSubmissions(context.Background(), &auth.UpdateUserSubmissionsRequest{
		Token:      initialToken,
		Queid:      "661fa6b218e230a331a2b428",
		Difficulty: "hard",
	})

	if err != nil {
		panic(err)
	}

	println(res2.User.GetEmail(), res2.String(), "htisfsdf")
}
