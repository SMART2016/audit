package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	auditlog "github.com/sirupsen/logrus"
	"log"
	"strings"
)

type EventLogNormalizer interface {
	RegisterLogPatterns(logtype string, pattern string)
	NormalizeLog(logtype string, logMsg string) string
}

type LogEventStoreClient interface {
	SubmitQuery(query string) ([]byte, error)
	PushLogEvents(logMsg string)
}

func startKafkaConsumer(logNormalizer EventLogNormalizer, eventStore LogEventStoreClient) {
	brokers := []string{"kafka:9092"}
	topic := "log_events_topic"
	logsChan := consumeKafkaMessages(brokers, topic)
	for msgMap := range logsChan {
		// Normalize your log message
		normalizedLog := logNormalizer.NormalizeLog(msgMap["msgType"].(string), string(msgMap["value"].([]byte)))
		if !strings.EqualFold(normalizedLog, "{}") {
			auditlog.Info(normalizedLog)
			eventStore.PushLogEvents(normalizedLog)
			//getNewElasticsearchClient().PushLogEvents(normalizedLog)
		}
	}
}

func consumeKafkaMessages(brokers []string, topic string) <-chan map[string]interface{} {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     topic,
		GroupID:   "audit-service",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(kafka.LastOffset)

	outChan := make(chan map[string]interface{})

	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("error while receiving message: %s\n", err.Error())
				continue
			}
			msgMap := map[string]interface{}{}
			msgMap["value"] = m.Value
			msgMap["msgType"] = string(m.Key)
			outChan <- msgMap
		}
	}()

	return outChan
}
