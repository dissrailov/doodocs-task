package main

import (
	"doodocs-task/config"
	"doodocs-task/internal/handlers"
	"doodocs-task/server"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	srv := server.NewServer(&cfg.HTTPServer, handlers.InitRoutes())
	log.Println("Starting server")
	if err := srv.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
