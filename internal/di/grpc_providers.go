package di

import (
	"context"
	"net"

	"go.uber.org/fx"
	"github.com/go-kit/log"
	"google.golang.org/grpc"

	helloDomain "engidone-auth/internal/hello/domain"
	helloEndpoints "engidone-auth/internal/hello/endpoints"
	helloPb "engidone-auth/internal/hello/proto"
	helloTransport "engidone-auth/internal/hello/transport"

	signinDomain "engidone-auth/internal/signin/domain"
	signinEndpoints "engidone-auth/internal/signin/endpoints"
	pb "engidone-auth/internal/signin/proto"
	signinTransport "engidone-auth/internal/signin/transport"
)

// GRPCModule provides all gRPC transport dependencies
var GRPCModule = fx.Options(
	fx.Provide(
		NewGRPCServer,
		NewHelloEndpoints,
		NewSigninEndpoints,
		NewHelloGRPCServer,
		NewSigninGRPCServer,
		NewTCPListener,
	),
	fx.Invoke(RegisterGRPCServices),
)

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// NewHelloEndpoints creates hello service endpoints
func NewHelloEndpoints(
	helloUC helloDomain.HelloUseCase,
	logger log.Logger,
) helloEndpoints.Set {
	return helloEndpoints.NewSet(helloUC)
}

// NewSigninEndpoints creates signin service endpoints
func NewSigninEndpoints(
	signinUC signinDomain.SigninUseCase,
	validateUC signinDomain.ValidateTokenUseCase,
	refreshUC signinDomain.RefreshTokenUseCase,
	getUserUC signinDomain.GetUserUseCase,
	logger log.Logger,
) signinEndpoints.Set {
	return signinEndpoints.NewSet(signinUC, validateUC, refreshUC, getUserUC, logger)
}

// NewHelloGRPCServer creates a hello service gRPC server
func NewHelloGRPCServer(endpoints helloEndpoints.Set) helloPb.HelloServiceServer {
	return helloTransport.NewGRPCServer(endpoints)
}

// NewSigninGRPCServer creates a signin service gRPC server
func NewSigninGRPCServer(endpoints signinEndpoints.Set) pb.SigninServiceServer {
	return signinTransport.NewGRPCServer(endpoints)
}

// NewTCPListener creates a TCP listener for the gRPC server
func NewTCPListener(config *AppConfig) (net.Listener, error) {
	address := ":" + config.ServerPort
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return lis, nil
}

// RegisterGRPCServices registers all gRPC services with the server and starts the server
func RegisterGRPCServices(
	lc fx.Lifecycle,
	grpcServer *grpc.Server,
	helloGRPCServer helloPb.HelloServiceServer,
	signinGRPCServer pb.SigninServiceServer,
	listener net.Listener,
	logger log.Logger,
	config *AppConfig,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Register gRPC services
			pb.RegisterSigninServiceServer(grpcServer, signinGRPCServer)
			helloPb.RegisterHelloServiceServer(grpcServer, helloGRPCServer)

			// Log startup information
			logger.Log("msg", "=== Engidone Auth Service ===")
			logger.Log("transport", "gRPC", "addr", ":"+config.ServerPort)
			logger.Log("msg", "Servidor iniciado en :"+config.ServerPort)
			logger.Log("msg", "Servicios disponibles:")
			logger.Log("msg", "  - Signin Service")
			logger.Log("msg", "  - Hello Service")
			logger.Log("msg", "")
			logger.Log("msg", "=== Usuarios disponibles para testing ===")
			logger.Log("msg", "Username: admin, Password: password123")
			logger.Log("msg", "Username: testuser, Password: test123")
			logger.Log("msg", "Username: john, Password: john123")
			logger.Log("msg", "=====================================")

			// Start gRPC server in a goroutine
			go func() {
				if err := grpcServer.Serve(listener); err != nil {
					logger.Log("transport", "gRPC", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Log("msg", "Stopping gRPC server")
			grpcServer.GracefulStop()
			return nil
		},
	})
}