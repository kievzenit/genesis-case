package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kievzenit/genesis-case/internal/api"
	"github.com/kievzenit/genesis-case/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := routes.RegisterRoutes()

	server := &http.Server{
		Addr:    cfg.ServerConfig.Address + ":" + fmt.Sprint(cfg.ServerConfig.Port),
		Handler: router,
	}

	go func() {
		log.Printf("starting server on %s:%d", cfg.ServerConfig.Address, cfg.ServerConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Kill)
	<-quitChan
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server shut down gracefully")
	log.Println("exiting")
}
