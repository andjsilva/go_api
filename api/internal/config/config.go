package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string
	RabbitMQDSN string
	QueueName   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente padrão")
	}

	return &Config{
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:123456@localhost:5432/postgres"),
		RabbitMQDSN: getEnv("RABBITMQ_DSN", "amqp://guest:guest@localhost:5672/"),
		QueueName:   getEnv("QUEUE_NAME", "product_queue"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
