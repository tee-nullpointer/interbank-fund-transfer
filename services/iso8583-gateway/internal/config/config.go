package config

import (
	"iso8583-gateway/pkg/util"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger      *LoggerConfig
	Server      *ServerConfig
	Kafka       *KafkaConfig
	Application *ApplicationConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type LoggerConfig struct {
	Level  string
	Format string
}

type KafkaConfig struct {
	Brokers []string
	Retry   int
	Timeout time.Duration
}

type ApplicationConfig struct {
	WorkerPerConnection int
	InboundRequestTopic string
	ServiceID           string
}

func Init() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	var serviceID string
	id, err := uuid.NewUUID()
	if err != nil {
		serviceID = util.RandomString(36)
	} else {
		serviceID = id.String()
	}
	return &Config{
		Logger: &LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Server: &ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "11111"),
		},
		Kafka: &KafkaConfig{
			Brokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}, ","),
			Retry:   getEnvAsInt("KAFKA_RETRY", 3),
			Timeout: getEnvAsDuration("KAFKA_TIMEOUT", 5*time.Second),
		},
		Application: &ApplicationConfig{
			WorkerPerConnection: getEnvAsInt("APP_WORKER_PER_CONNECTION", 10),
			InboundRequestTopic: getEnv("APP_INBOUND_REQUEST_TOPIC", "transfer.inbound.request"),
			ServiceID:           serviceID,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, sep)
	}
	return defaultValue
}
