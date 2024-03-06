package main

import (
	server "github.com/WildEgor/gAuth/internal"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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

	srv, _ := server.NewServer()

	log.Println("[Main] Starting gRPC server")
	err := srv.GRPC.Init()
	if err != nil {
		log.Panic(err)
	}

	log.Println("[Main] Connect to Mongo and Redis")
	srv.Mongo.Connect()
	srv.Redis.Connect()

	addr := ":" + srv.AppConfig.Port
	if err := srv.App.Listen(addr); err != nil {
		log.Panicf("[CRIT] Unable to start server. Reason: %v", err)
	}
	log.Printf("Server is listening on PORT: %s", srv.AppConfig.Port)

	<-done

	log.Println("[Main] Stopping consumer")
	if err := srv.GRPC.Stop; err != nil {
		log.Panicf("[CRIT] Unable to start gRPC server. Reason: %v", err)
	}
	srv.Redis.Disconnect()
	srv.Mongo.Disconnect()

	if srv.Notifier.Close() != nil {
		log.Panicf("[CRIT] Unable to close notifier. Reason: %v", err)
	}

	if err := srv.App.Shutdown(); err != nil {
		log.Panicf("[CRIT] Unable to shutdown server. Reason: %v", err)
	}
}
