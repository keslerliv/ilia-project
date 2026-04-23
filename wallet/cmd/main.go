package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/keslerliv/ilia-project/wallet/config"
	"github.com/keslerliv/ilia-project/wallet/internal/models"
	"github.com/keslerliv/ilia-project/wallet/internal/routes"
	"github.com/keslerliv/ilia-project/wallet/internal/services/kafka"
	"github.com/keslerliv/ilia-project/wallet/pkg/db"
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
	db.MakeMigration(conn)

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

	consumer, err := user.ConsumePartition(config.Env.KafkaTopic, 0, sarama.OffsetNewest)
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
				ProcessMessage(msg)
			case <-sigchan:
				fmt.Println("Shutting down...")
				close(doneCh)
				return
			}
		}
	}()
}

func ProcessMessage(msg *sarama.ConsumerMessage) {
	var payload map[string]interface{}
	message := strings.NewReader(string(msg.Value))

	err := json.NewDecoder(message).Decode(&payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create new user wallet
	if payload["action"] == "new_user" {
		walletId, err := models.CreateWallet(int(payload["user_id"].(float64)))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("New wallet %v created for user %v", walletId, payload["user_id"])
		return
	}
}
