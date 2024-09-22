package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/Nelwhix/kafka-notify/cmd/producer/handlers"
	"github.com/Nelwhix/kafka-notify/pkg"
	"github.com/Nelwhix/kafka-notify/pkg/middlewares"
	"github.com/Nelwhix/kafka-notify/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

const (
	ProducerPort    = ":8080"
	KafkaServerAddr = "localhost:9092"
)

var validate *validator.Validate
var decoder = schema.NewDecoder()

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

	logger, err := pkg.CreateNewLogger()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	model := models.Model{
		Conn: conn,
	}

	handler := handlers.Handler{
		Producer:  producer,
		Model:     model,
		Logger:    logger,
		Decoder:   decoder,
		Validator: validate,
	}

	r := http.NewServeMux()
	r.HandleFunc("POST /send", handler.SendMessage)

	fmt.Printf("Kafka Producer started at http://localhost:%s\n", ProducerPort)

	err = http.ListenAndServe(ProducerPort, middlewares.ContentTypeMiddleware(r))

	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
