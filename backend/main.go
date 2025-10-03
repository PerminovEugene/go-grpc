package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-grpc-backend/internal/server"
)

func main() {
	port := getEnv("GRPC_PORT", "50051")

	server, err := server.NewAnalyticsServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		server.Stop()
		os.Exit(0)
	}()

	log.Printf("Starting server on port %s", port)
	if err := server.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
