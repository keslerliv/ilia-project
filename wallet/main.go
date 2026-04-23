// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/IBM/sarama"
// )

// type Order struct {
// 	CustomerName string `json:"customer_name"`
// 	CoffeeType   string `json:"coffee_type"`
// }

// func main() {
// 	http.HandleFunc("/place-order", placeOrder)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func Connectwallet(brokers []string) (sarama.Syncwallet, error) {
// 	config := sarama.NewConfig()
// 	config.wallet.Return.Successes = true
// 	config.wallet.RequiredAcks = sarama.WaitForAll
// 	config.wallet.Retry.Max = 5

// 	wallet, err := sarama.NewSyncProducer(brokers, config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return wallet, nil

// }

// func PushOrderToQueue(topic string, orderBytes []byte) error {
// 	brokers := []string{"kafka:9092"}
// 	wallet, err := Connectwallet(brokers)
// 	if err != nil {
// 		return err
// 	}
// 	defer wallet.Close()

// 	// Create a new message to send to the Kafka topic
// 	msg := &sarama.walletMessage{
// 		Topic: topic,
// 		Value: sarama.StringEncoder(orderBytes),
// 	}

// 	// Send message to Kafka topic
// 	partition, offset, err := wallet.SendMessage(msg)
// 	if err != nil {
// 		return err
// 	}

// 	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)

// 	return nil
// }

// func placeOrder(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// 1. Parse request body into order
// 	var order Order
// 	err := json.NewDecoder(r.Body).Decode(&order)
// 	if err != nil {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}

// 	// 2. Conver body into bytes
// 	orderBytes, err := json.Marshal(order)
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	// 3. Send order to Kafka topic
// 	err = PushOrderToQueue("coffee-orders", orderBytes)
// 	if err != nil {
// 		http.Error(w, "Failed to place order", http.StatusInternalServerError)
// 		return
// 	}

// 	// 4. Respond to client
// 	response := map[string]interface{}{
// 		"message": "Order placed successfully",
// 		"order":   order,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		return
// 	}
// }
