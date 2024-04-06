package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

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

func createElasticsearchClient(addresses string) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: strings.Split(addresses, ","),
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return es
}
