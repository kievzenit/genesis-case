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

	"github.com/go-co-op/gocron/v2"
	_ "github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kievzenit/genesis-case/internal/api"
	"github.com/kievzenit/genesis-case/internal/config"
	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/jobs"
	"github.com/kievzenit/genesis-case/internal/models"
	"github.com/kievzenit/genesis-case/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create scheduler: %v", err)
	}

	weatherService := services.NewWeatherService(cfg.WeatherServiceConfig)

	emailService := services.NewEmailService(cfg.BaseURL, cfg.EmailServiceConfig)

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

	txManager := database.NewTransactionManager(sqlCon)

	router := routes.RegisterRoutes(
		weatherService,
		emailService,
		sqlCon,
		txManager,
	)

	server := &http.Server{
		Addr:         cfg.ServerConfig.Address + ":" + fmt.Sprint(cfg.ServerConfig.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeout) * time.Second,
	}

	sendConfirmationEmailJob := jobs.NewSendConfirmationEmailJob(
		emailService,
		txManager,
	)

	_, err = scheduler.NewJob(
		gocron.DurationJob(time.Duration(cfg.JobsConfig.EmailConfirmationInterval)*time.Minute),
		gocron.NewTask(sendConfirmationEmailJob.Run),
	)
	if err != nil {
		log.Fatalf("failed to create send confirmation email job: %v", err)
	}

	sendWeatherReportJob := jobs.NewSendWeatherReportJob(
		weatherService,
		emailService,
		sqlCon,
	)

	_, err = scheduler.NewJob(
		gocron.DurationJob(time.Duration(1)*time.Hour),
		gocron.NewTask(sendWeatherReportJob.Run, models.Hourly),
	)
	if err != nil {
		log.Fatalf("failed to create send hourly weather report job: %v", err)
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(12, 0, 0))),
		gocron.NewTask(sendWeatherReportJob.Run, models.Daily),
	)
	if err != nil {
		log.Fatalf("failed to create send daily weather report job: %v", err)
	}

	go func() {
		log.Printf("starting server on %s:%d", cfg.ServerConfig.Address, cfg.ServerConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	scheduler.Start()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, os.Kill)
	<-quitChan
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	if err := scheduler.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown scheduler: %v", err)
	}

	log.Println("server shut down gracefully")
	log.Println("exiting")
}
