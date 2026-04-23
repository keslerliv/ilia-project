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

	JWTSecret string
}

var Env AppConfig

func LoadConfig() *AppConfig {
	// Server
	Env.Port = "8081"

	// Database
	Env.DBHost = "postgres"
	Env.DBPort = "5432"
	Env.DBUser = "postgres"
	Env.DBPassword = "postgres"
	Env.DBName = "ilia"
	Env.DBSSLMode = "disable"

	// Kafka
	Env.KafkaBrokers = []string{"localhost:9092"}
	Env.KafkaTopic = "transaction"

	// JWT
	Env.JWTSecret = "mysecretkey"

	return &Env
}
