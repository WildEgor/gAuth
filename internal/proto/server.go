package proto

import (
	"fmt"
	"github.com/WildEgor/gAuth/internal/configs"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	server    *grpc.Server
	appConfig *configs.AppConfig
}

func NewGRPCServer(
	appConfig *configs.AppConfig,
) *GRPCServer {
	return &GRPCServer{
		appConfig: appConfig,
	}
}

func (s *GRPCServer) Init() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.appConfig.RPCPort))
	if err != nil {
		log.Fatal("cannot listen port - ", err)
	}

	serv := grpc.NewServer()

	RegisterCentrifugoProxyServer(serv, &ProxyService{})
	RegisterAuthServiceServer(serv, &AuthService{})

	s.server = serv

	go func() {
		// Start gRPC server
		if err := serv.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (s *GRPCServer) Stop() {
	s.server.Stop()
}
