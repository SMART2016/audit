package main

import (
	"encoding/json"
	"github.com/vjeantet/grok"
	"log"
)

func normalizeLog(logMsg string) string {
	g, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	if err != nil {
		log.Fatalf("Failed to create Grok instance: %s", err)
	}

	pattern := "%{COMBINEDAPACHELOG}" // Example pattern for Apache logs
	values, err := g.Parse(pattern, logMsg)
	if err != nil {
		log.Printf("Failed to parse log message: %s", err)
		return ""
	}

	// Assuming you convert values to a JSON string or another structured format
	normalizedLog, err := json.Marshal(values)
	if err != nil {
		log.Printf("Failed to marshal values: %s", err)
		return ""
	}

	return string(normalizedLog)
}
