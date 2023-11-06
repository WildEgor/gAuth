package main

import (
	"fmt"
	server "github.com/WildEgor/gAuth/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	srv, _ := server.NewServer()
	log.Fatal(srv.Listen(fmt.Sprintf(":%v", "8888")))
}
