package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultDBHost     = "localhost"
	defaultDBPort     = 5432
	defaultDBUser     = "postgres"
	defaultDBPassword = "password"
	defaultDBName     = "postgres"
	defaultDBTable    = "test_data"

	defaultKafkaBrokers = "localhost:9092"
	defaultKafkaTopic   = "test-topic"

	defaultHTTPPort         = 8080
	defaultPollDelay        = 5 * time.Second
	defaultHTTPStartupDelay = 5 * time.Second
)

// PostgresConfig holds PostgreSQL connection settings.
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Table    string
}

// KafkaConfig holds Kafka connection settings.
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// Config holds service configuration.
type Config struct {
	Postgres         PostgresConfig
	Kafka            KafkaConfig
	HTTPPort         int
	HTTPStartupDelay time.Duration
	PollInterval     time.Duration
}

// FromEnv loads configuration from environment variables.
func FromEnv() Config {
	brokers := strings.Split(getEnv("KAFKA_BROKERS", defaultKafkaBrokers), ",")

	return Config{
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", defaultDBHost),
			Port:     getEnvInt("POSTGRES_PORT", defaultDBPort),
			User:     getEnv("POSTGRES_USER", defaultDBUser),
			Password: getEnv("POSTGRES_PASSWORD", defaultDBPassword),
			DBName:   getEnv("POSTGRES_DB", defaultDBName),
			Table:    getEnv("POSTGRES_TABLE", defaultDBTable),
		},
		Kafka: KafkaConfig{
			Brokers: brokers,
			Topic:   getEnv("KAFKA_TOPIC", defaultKafkaTopic),
		},
		HTTPPort:         getEnvInt("HTTP_PORT", defaultHTTPPort),
		HTTPStartupDelay: getEnvDuration("HTTP_STARTUP_DELAY", defaultHTTPStartupDelay),
		PollInterval:     getEnvDuration("POLL_INTERVAL", defaultPollDelay),
	}
}

// HTTPAddr returns the HTTP bind address.
func (c Config) HTTPAddr() string {
	return fmt.Sprintf(":%d", c.HTTPPort)
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
