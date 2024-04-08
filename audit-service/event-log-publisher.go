package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

const (
	kafkaTopic   = "log_events_topic"
	kafkaBroker  = "kafka:9092"
	kafkaTimeout = 10 * time.Second
)

func publishEventLogs(logType string, logMsg string) {
	fmt.Println("Sending Log Event = ", logMsg)
	// Initialize a Kafka writer with the broker and topic.
	w := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(fmt.Sprintf("%s", logType)),
			Value: []byte(logMsg),
		},
	)

	if err != nil {
		log.Fatalf("failed to write messages: %v", err)
	}
	log.Printf("sent: %s", logMsg)

}
