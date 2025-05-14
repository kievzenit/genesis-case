package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fmt"

	"database/sql"
	_ "github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kievzenit/genesis-case/internal/api"
	"github.com/kievzenit/genesis-case/internal/config"
	"github.com/kievzenit/genesis-case/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	weatherService := services.NewWeatherService(cfg.WeatherServiceConfig.ApiKey, cfg.WeatherServiceConfig.HttpTimeout)

	sqlCon, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Username,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.DatabaseName,
	))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer sqlCon.Close()
	if err := sqlCon.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	driver, err := postgres.WithInstance(sqlCon, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create database driver: %v", err)
	}
	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	if cfg.DatabaseConfig.ApplyMigrations {
		log.Println("applying migrations...")
		if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to apply migrations: %v", err)
		}
	}

	router := routes.RegisterRoutes(weatherService)

	server := &http.Server{
		Addr:         cfg.ServerConfig.Address + ":" + fmt.Sprint(cfg.ServerConfig.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeout) * time.Second,
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
