package main

import (
	"fmt"
	server "github.com/WildEgor/gAuth/internal"
	"github.com/WildEgor/gAuth/internal/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	srv, _ := server.NewServer()

	grpc := proto.GRPCServer{}

	init, err := grpc.Init()
	if err != nil {
		return
	}
	defer init.Stop()

	log.Fatal(srv.Listen(fmt.Sprintf(":%v", "8888")))
}
