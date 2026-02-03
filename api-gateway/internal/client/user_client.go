package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/adityakw90/service-user-proto/gen/go/user"
	authpb "github.com/adityakw90/service-user-proto/gen/go/auth"
)

// UserClient wraps the gRPC clients for user and auth services
type UserClient struct {
	userClient userpb.UserServiceClient
	authClient authpb.AuthServiceClient
	conn       *grpc.ClientConn
}

// Config holds the configuration for the user client
type Config struct {
	ServiceAddress string
	DialTimeout    time.Duration
}

// DefaultConfig returns the default configuration for the user client
func DefaultConfig() Config {
	return Config{
		ServiceAddress: "localhost:50051",
		DialTimeout:    5 * time.Second,
	}
}

// NewUserClient creates a new user client with the given configuration
func NewUserClient(cfg Config) (*UserClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		cfg.ServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024), // 4MB
			grpc.MaxCallSendMsgSize(4*1024*1024), // 4MB
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial user service: %w", err)
	}

	return &UserClient{
		userClient: userpb.NewUserServiceClient(conn),
		authClient: authpb.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil
}

// Close closes the gRPC connection
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// UserService returns the user service client
func (c *UserClient) UserService() userpb.UserServiceClient {
	return c.userClient
}

// AuthService returns the auth service client
func (c *UserClient) AuthService() authpb.AuthServiceClient {
	return c.authClient
}
