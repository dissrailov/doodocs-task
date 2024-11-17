package main

import (
	"doodocs-task/config"
	"doodocs-task/internal/handlers"
	"doodocs-task/internal/service"
	"doodocs-task/server"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	service := service.NewService(&cfg.SMTP)
	handler := handlers.NewHandler(service)

	router := handlers.InitRoutes(handler)
	srv := server.NewServer(&cfg.HTTPServer, router)

	log.Printf("server is running on: %s", cfg.HTTPServer.Addr)

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
