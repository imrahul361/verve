package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var kafkaWriter *kafka.Writer

func sendToKafka(count int) error {
	payload, err := json.Marshal(map[string]int{"count": count})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("Timestamp::%d", time.Now().Unix())),
		Value: payload,
	}

	return kafkaWriter.WriteMessages(context.Background(), msg)
}

// CountLoggerIntoKafka : Log the count of unique IDs every minute into Kafka
func CountLoggerIntoKafka() {
	err := sendToKafka(getCount())
	if err != nil {
		log.Printf("Failed to send to Kafka: %v", err)
	}
	// re-initialise it for the next minute after logging
	resetCount()
}

// InitKafka : Initialize Kafka writer
func InitKafka() {

	kafkaWriter = &kafka.Writer{
		Addr:       kafka.TCP("127.0.0.1:9092"),
		Topic:      "unique_request_id_log",
		BatchSize:  100,
		Async:      true,
		Completion: nil,
	}
}
