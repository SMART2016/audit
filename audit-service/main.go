package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strconv"
	"time"

	"strings"
)

func main() {
	brokers := []string{"localhost:9093"}
	topic := "log_events_topic"
	elasticsearchAddresses := "http://localhost:9200"

	esClient := createElasticsearchClient(elasticsearchAddresses)
	logsChan := consumeKafkaMessages(brokers, topic)

	for logMsg := range logsChan {
		// Normalize your log message
		normalizedLog := normalizeLog(logMsg)
		fmt.Println("normal log", normalizedLog)
		// Push to Elasticsearch
		req := esapi.IndexRequest{
			Index:      "your-index-name",
			DocumentID: strconv.Itoa(time.Now().Nanosecond()), // Example ID, consider a better one
			Body:       strings.NewReader(normalizedLog),
			Refresh:    "true",
		}
		res, err := req.Do(context.Background(), esClient)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			log.Printf("[%s] Error indexing document ID=%d", res.Status(), req.DocumentID)
		} else {
			log.Printf("Document ID=%d indexed.", req.DocumentID)
		}
	}
}
