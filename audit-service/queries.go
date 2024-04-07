package main

import (
	"encoding/json"
	"fmt"
)

func generateAndFilterQuery(originalQueryJSON string, filterMap map[string]interface{}) (string, error) {
	// Original query JSON
	if originalQueryJSON == "" {
		originalQueryJSON = `{"query": {"range": {"time": {"gte": "now-48h","lte": "now"}}}}`
	}
	fmt.Println("filterMap =", filterMap)
	// Unmarshal the original query into a map
	var originalQueryMap map[string]interface{}
	err := json.Unmarshal([]byte(originalQueryJSON), &originalQueryMap)
	if err != nil {
		panic(err) // Handle error appropriately
	}

	// Define the additional conditions to be added
	filterLst := []interface{}{}
	for k, v := range filterMap {
		filterLst = append(filterLst, map[string]interface{}{
			"match": map[string]interface{}{
				k: v,
			},
		})
	}

	// Insert the additional conditions into the original query map
	// Check if a bool query needs to be constructed
	if query, exists := originalQueryMap["query"].(map[string]interface{}); exists {
		// Construct the bool query with the original range query and additional conditions
		boolQuery := map[string]interface{}{
			"bool": map[string]interface{}{
				"must": append(filterLst, query),
			},
		}

		// Replace the original query with the new bool query
		originalQueryMap["query"] = boolQuery
	} else {
		panic("Unexpected structure of original query") // Handle appropriately
	}

	// Marshal the modified query map back to JSON to see the result
	modifiedQueryJSON, err := json.MarshalIndent(originalQueryMap, "", "  ")
	if err != nil {
		panic(err) // Handle error appropriately
	}
	return string(modifiedQueryJSON), nil
}
