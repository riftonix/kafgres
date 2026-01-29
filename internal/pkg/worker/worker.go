package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kafgres/internal/pkg/db"
	"kafgres/internal/pkg/health"
	"kafgres/internal/pkg/kafka"

	kafkaGo "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Data represents one DB record.
type Data struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Start runs a periodic ticker.
func Start(interval time.Duration, fn func()) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			fn()
		}
	}()
}

// ProcessCycle reads from DB, writes to Kafka, and updates health.
func ProcessCycle(database db.Interface, writer kafka.Writer, state *health.State, table string) {
	data, err := readFromDB(database, table)
	if err != nil {
		logrus.Errorf("Error reading from DB: %v", err)
		state.SetHealthy(false)
		return
	}

	if len(data) == 0 {
		logrus.Info("No data found in database")
		state.SetHealthy(false)
		return
	}

	if err := writeToKafka(writer, data); err != nil {
		logrus.Errorf("Error writing to Kafka: %v", err)
		state.SetHealthy(false)
		return
	}

	state.SetHealthy(true)
}

// WaitForShutdown blocks on SIGINT/SIGTERM and runs cleanup.
func WaitForShutdown(ctx context.Context, cleanup func() error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logrus.WithField("signal", sig.String()).Info("Shutdown signal received, shutting down gracefully...")
	if err := cleanup(); err != nil {
		logrus.Fatalf("Shutdown error: %v", err)
	}
}

func readFromDB(database db.Interface, table string) ([]Data, error) {
	rows, err := database.Query(fmt.Sprintf("SELECT id, name FROM %s LIMIT 100", table))
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	var data []Data
	for rows.Next() {
		var d Data
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		data = append(data, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	logrus.Infof("Successfully read %d records from database", len(data))
	return data, nil
}

func writeToKafka(writer kafka.Writer, data []Data) error {
	messages := make([]kafkaGo.Message, len(data))
	for i, d := range data {
		jsonData, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}

		messages[i] = kafkaGo.Message{Value: jsonData}
	}

	if err := writer.WriteMessages(context.Background(), messages...); err != nil {
		return fmt.Errorf("failed to write to kafka: %w", err)
	}

	logrus.Infof("Successfully wrote %d records to kafka", len(data))
	return nil
}
