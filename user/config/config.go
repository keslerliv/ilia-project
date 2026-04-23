package config

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
}

var Env AppConfig

func LoadConfig() *AppConfig {
	// Server
	Env.Port = "8081"

	// Database
	Env.DBHost = "localhost"
	Env.DBPort = "5432"
	Env.DBUser = "postgres"
	Env.DBPassword = "password"
	Env.DBName = "ilia"
	Env.DBSSLMode = "disable"

	// Kafka
	Env.KafkaBrokers = []string{"localhost:9092"}
	Env.KafkaTopic = "transaction"

	return &Env
}
