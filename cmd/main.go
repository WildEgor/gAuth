package main

import (
	"errors"
	"fmt"
	server "github.com/WildEgor/gAuth/internal"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var srv *fiber.App

func main() {
	Start()
	Shutdown()
}

func Start() {
	srv, _ = server.NewServer()
	go func() {
		if err := srv.Listen(fmt.Sprintf(":%v", "8888")); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func Shutdown() {
	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println()
		log.Println(sig)
		done <- true
	}()

	log.Println("[Main] Awaiting signal")
	<-done
	log.Println("[Main] Stopping consumer")
}
