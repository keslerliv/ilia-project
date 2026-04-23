package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/keslerliv/ilia-project/user/config"
)

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	user, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func PushOrderToQueue(orderBytes []byte) error {
	user, err := ConnectProducer(config.Env.KafkaBrokers)
	if err != nil {
		return err
	}
	defer user.Close()

	// Create a new message to send to the Kafka topic
	msg := &sarama.ProducerMessage{
		Topic: config.Env.KafkaTopic,
		Value: sarama.StringEncoder(orderBytes),
	}

	// Send message to Kafka topic
	partition, offset, err := user.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	return nil
}
