package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	USER_SERVICE_LOG_TYPE    = "user-service"
	USER_SERVICE_LOG_PATTERN = `CurrentUser: %{WORD:currentUser}, System: %{USERNAME:system}, Action: %{WORD:action}, IP: \[%{IPV6:ip}\]:%{NUMBER:port}, Agent: %{GREEDYDATA:agent}, Time: %{TIMESTAMP_ISO8601:time}`
)

func init() {
	file, err := os.OpenFile("service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}

	log.SetOutput(file)
	// Set log level
	log.SetLevel(log.InfoLevel)

	// Use JSON formatter
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	brokers := []string{"localhost:9093"}
	topic := "log_events_topic"
	//elasticsearchAddresses := "http://localhost:9200"
	logNormalizer := LogNormalizer{make(map[string]string)}
	logNormalizer.registerLogPatterns(USER_SERVICE_LOG_TYPE, USER_SERVICE_LOG_PATTERN)
	//esClient := createElasticsearchClient(elasticsearchAddresses)
	logsChan := consumeKafkaMessages(brokers, topic)

	for logMsg := range logsChan {
		// Normalize your log message
		normalizedLog := logNormalizer.normalizeLog(USER_SERVICE_LOG_TYPE, logMsg)
		log.Info(normalizedLog)
		// Push to Elasticsearch
		//req := esapi.IndexRequest{
		//	Index:      "index-log-events",
		//	DocumentID: strconv.Itoa(time.Now().Nanosecond()), // Example ID, consider a better one
		//	Body:       strings.NewReader(normalizedLog),
		//	Refresh:    "true",
		//}
		//res, err := req.Do(context.Background(), esClient)
		//if err != nil {
		//	log.Fatalf("Error getting response: %s", err)
		//}
		//defer res.Body.Close()
		//if res.IsError() {
		//	log.Printf("[%s] Error indexing document ID=%d", res.Status(), req.DocumentID)
		//} else {
		//	log.Printf("Document ID=%d indexed.", req.DocumentID)
		//}
	}
}
