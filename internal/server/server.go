package server

import (
	"engidoneauth/internal/app"
	"engidoneauth/internal/config"
	"engidoneauth/internal/credentials"
	"engidoneauth/internal/db"
	"engidoneauth/internal/greet"
	"engidoneauth/internal/jwt"
	pb "engidoneauth/internal/proto"
	"engidoneauth/internal/signin"
	"engidoneauth/internal/users"
	"engidoneauth/log"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	appConfig *config.AppConfig
	certs     jwt.Certs
	users     []users.User
}

func NewGRPCServer(appConfig *config.AppConfig, certs jwt.Certs, users []users.User) *GRPCServer {

	s := &GRPCServer{
		appConfig,
		certs,
		users,
	}

	grpcServer := grpc.NewServer()
	dbModule := s.dbModule()
	greeModule := greetModule()
	signinModule := s.signInModule(dbModule)
	jwtModule := s.jwtModule(dbModule)
	application := app.NewUseCase(greeModule, signinModule, jwtModule)

	pb.RegisterAuthServiceServer(grpcServer, application)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.Application.Server.Port))
	if err != nil {
		panic(err)
	}

	log.Infof("gRPC server listening on :%s", appConfig.Application.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return s
}

func greetModule() *greet.UseCase {
	return greet.NewUseCase(greet.NewGreetingRepository())
}

func (s *GRPCServer) jwtModule(dbModule *db.Queries) *jwt.UseCase {
	repository := jwt.NewSQLRepository(dbModule)
	return jwt.NewUseCase(60*time.Minute, s.certs, repository)
}

func credentialsModule(dbModule *db.Queries) *credentials.UseCase {
	repository := credentials.NewSQLRepository(dbModule)
	return credentials.NewUseCase(repository)
}

func (s *GRPCServer) usersModule() *users.UseCase {
	repository := users.NewRPCUserServiceRepository(s.users)
	return users.NewUseCase(repository)
}

func (s *GRPCServer) signInModule(dbModule *db.Queries) *signin.UseCase {
	return signin.NewUseCase(s.jwtModule(dbModule), credentialsModule(dbModule), s.usersModule())
}

func (s *GRPCServer) dbModule() *db.Queries {
	dbconn, _ := db.NewDBConnection(s.appConfig)
	return db.New(dbconn)
}
