package app

import (
	"context"
	"net/http"
	"time"

	"kafgres/internal/pkg/config"
	"kafgres/internal/pkg/db"
	"kafgres/internal/pkg/health"
	"kafgres/internal/pkg/kafka"
	"kafgres/internal/pkg/worker"

	"github.com/sirupsen/logrus"
)

// Init bootstraps the service.
func Init() {
	logrus.Info("Starting kafgres service...")

	cfg := config.FromEnv()

	database, err := db.Connect(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer database.Close()

	kafkaWriter, err := kafka.Connect(cfg.Kafka)
	if err != nil {
		logrus.Fatalf("Failed to initialize Kafka connection: %v", err)
	}
	defer kafkaWriter.Close()

	healthState := health.NewState()

	worker.Start(cfg.PollInterval, func() {
		worker.ProcessCycle(database, kafkaWriter, healthState, cfg.Postgres.Table)
	})

	server := &http.Server{
		Addr: cfg.HTTPAddr(),
	}

	http.HandleFunc("/health", health.Handler(healthState))

	go func() {
		logrus.Infof("Starting HTTP server on %s", cfg.HTTPAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("HTTP server error: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	worker.WaitForShutdown(ctx, func() error {
		return server.Shutdown(ctx)
	})
}
