package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/keslerliv/wallet/config"
	"github.com/keslerliv/wallet/internal/routes"
	"github.com/keslerliv/wallet/internal/services/kafka"
	"github.com/keslerliv/wallet/pkg/db"
)

func main() {
	config.LoadConfig()

	// Start database connection
	conn, err := db.OpenConnection()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Make initial database migration
	err = db.MakeMigration(conn)
	if err != nil {
		panic(err)
	}

	// Call routine to listen for Kafka messages
	// This function will listen for new users and create the wallet
	ListenMessages()

	r := routes.LoadRoutes()

	// Start server
	http.ListenAndServe(fmt.Sprintf(":%s", config.Env.Port), r)
}

func ListenMessages() {
	// Create new consumer
	user, err := kafka.ConnectConsumer(config.Env.KafkaBrokers)
	if err != nil {
		panic(err)
	}

	consumer, err := user.ConsumePartition(config.Env.KafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to consume user update messages
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Printf("Error consuming message: %v\n", err)
				continue
			case msg := <-consumer.Messages():
				fmt.Printf("Received message: %s\n", msg.Value)
				order := string(msg.Value)
				fmt.Printf("Brewing order: %s\n", order)
			case <-sigchan:
				fmt.Println("Shutting down...")
				close(doneCh)
				return
			}
		}
	}()
}
