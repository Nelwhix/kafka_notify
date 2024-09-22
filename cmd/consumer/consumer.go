package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Nelwhix/kafka-notify/cmd/consumer/handlers"
	"github.com/Nelwhix/kafka-notify/pkg/middlewares"
	"github.com/Nelwhix/kafka-notify/pkg/models"
	"log"
	"net/http"
)

const (
	Group           = "notifications-group"
	KafkaServerAddr = "localhost:9092"
	Port            = ":8081"
)

type Consumer struct {
	store *models.NotificationStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		userID := string(msg.Key)
		var notification models.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}

		consumer.store.Add(userID, notification)
		sess.MarkMessage(msg, "")
	}

	return nil
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup([]string{KafkaServerAddr}, Group, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return consumerGroup, nil
}

func setupConsumerGroup(ctx context.Context, store *models.NotificationStore) {
	consumerGroup, err := initializeConsumerGroup()
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		store: store,
	}

	for {
		err = consumerGroup.Consume(ctx, []string{"notifications"}, consumer)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func main() {
	store := &models.NotificationStore{
		Data: make(models.UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroup(ctx, store)
	defer cancel()

	handler := handlers.Handler{
		NotificationStore: store,
	}

	r := http.NewServeMux()
	r.HandleFunc("GET /notifications/{userId}", handler.GetNotifications)

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+"started at http://localhost%s\n", Group, Port)

	err := http.ListenAndServe(Port, middlewares.ContentTypeMiddleware(r))

	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
