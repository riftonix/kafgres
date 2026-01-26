package kafka

import (
	"context"
	"kafgres/internal/pkg/config"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Writer abstracts Kafka writing for testing.
type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Connect creates a Kafka writer.
func Connect(cfg config.KafkaConfig) (Writer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
	})

	logrus.Info("Successfully connected to Kafka")
	return writer, nil
}
