package config

import "os"

type AppConfig struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	KafkaBrokers []string
	KafkaTopic   string

	JWTSecret string
}

var Env AppConfig

func LoadConfig() *AppConfig {
	// Server
	Env.Port = os.Getenv("WALLET_APP_PORT")

	// Database
	Env.DBHost = os.Getenv("WALLET_DB_HOST")
	Env.DBPort = os.Getenv("WALLET_DB_PORT")
	Env.DBUser = os.Getenv("WALLET_DB_USER")
	Env.DBPassword = os.Getenv("WALLET_DB_PASSWORD")
	Env.DBName = os.Getenv("WALLET_DB_NAME")

	// Kafka
	Env.KafkaBrokers = []string{os.Getenv("KAFKA_BROKERS")}
	Env.KafkaTopic = os.Getenv("KAFKA_TOPIC")

	// JWT
	Env.JWTSecret = os.Getenv("JWT_SECRET")

	return &Env
}
