package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	auditlog "github.com/sirupsen/logrus"
	"log"
)

func startKafkaConsumer(logNormalizer LogNormalizer) {
	brokers := []string{"localhost:9093"}
	topic := "log_events_topic"
	logsChan := consumeKafkaMessages(brokers, topic)
	for logMsg := range logsChan {
		// Normalize your log message
		normalizedLog := logNormalizer.normalizeLog(USER_SERVICE_LOG_TYPE, logMsg)
		auditlog.Info(normalizedLog)
		getNewElasticsearchClient().pushLogEvents(normalizedLog)
	}
}

func consumeKafkaMessages(brokers []string, topic string) <-chan string {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		GroupID:   "audit-service",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(kafka.LastOffset)

	outChan := make(chan string)

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(context.Background())
			fmt.Println(string(m.Value))
			if err != nil {
				log.Printf("error while receiving message: %s\n", err.Error())
				continue
			}
			outChan <- string(m.Value)
		}
	}()

	return outChan
}
