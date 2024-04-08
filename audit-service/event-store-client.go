package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	ES_ADDRESS = "http://localhost:9200"
)

type EsClient struct {
	esClient *elasticsearch.Client
}

func getNewElasticsearchClient() *EsClient {
	cfg := elasticsearch.Config{
		Addresses: strings.Split(ES_ADDRESS, ","),
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return &EsClient{es}
}

func (es *EsClient) PushLogEvents(logMsg string) {
	// Push to Elasticsearch
	req := esapi.IndexRequest{
		Index:      "index-log-events",
		DocumentID: strconv.Itoa(time.Now().Nanosecond()), // Example ID, consider a better one
		Body:       strings.NewReader(logMsg),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es.esClient)
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

func (es *EsClient) SubmitQuery(query string) ([]byte, error) {
	fmt.Println("ES Query to Submit = ", query)
	req := esapi.SearchRequest{
		Index:  []string{"index-log-events*"}, // Adjust with your index pattern
		Body:   strings.NewReader(query),
		Pretty: true,
	}
	res, err := req.Do(context.Background(), es.esClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()
	var respBytes []byte
	if res.IsError() {
		var r map[string]interface{}
		json.NewDecoder(res.Body).Decode(&r)

	} else {

		// Parse the response body
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			return nil, err
		} else {
			srcLst := r["hits"].(map[string]interface{})["hits"].([]interface{})
			respLst := []map[string]interface{}{}
			for _, source := range srcLst {
				respLst = append(respLst, source.(map[string]interface{})["_source"].(map[string]interface{}))
			}
			respBytes, _ = json.Marshal(respLst)
			return respBytes, nil
		}
	}
	return respBytes, nil
}
