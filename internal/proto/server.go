package proto

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func Init() (*grpc.Server, error) {
	const port = ":8088" // TODO: use config
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("cannot listen port - ", err)
	}
	serv := grpc.NewServer()
	RegisterCentrifugoProxyServer(serv, &ProxyService{})
	RegisterAuthServiceServer(serv, &AuthService{})

	go func() {
		// Start gRPC server
		if err := serv.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return serv, nil
}
