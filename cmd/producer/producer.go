package main

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Nelwhix/kafka-notify/cmd/producer/handlers"
	"github.com/Nelwhix/kafka-notify/pkg"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

const (
	ProducerPort    = ":8080"
	KafkaServerAddr = "localhost:9092"
	KafkaTopic      = "notifications"
)

var ErrUserNotFoundInProducer = errors.New("user not found")

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddr}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}

	return producer, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()

	conn, err := pkg.CreateDbConn()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	handler := handlers.Handler{
		Producer: producer,
		Conn:     conn,
	}

	http.HandleFunc("POST /send", handler.SendMessage)

	fmt.Printf("Kafka Producer started at http://localhost:%s/\n", ProducerPort)

	err = http.ListenAndServe(ProducerPort, nil)

	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
