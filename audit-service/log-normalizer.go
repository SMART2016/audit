package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vjeantet/grok"
)

type LogNormalizer struct {
	patternMap map[string]string
}

func (l *LogNormalizer) registerLogPatterns(logtype string, pattern string) {
	l.patternMap[logtype] = pattern
}
func (l *LogNormalizer) normalizeLog(logtype string, logMsg string) string {
	fmt.Println("logMsg to normalize =", logMsg)
	g, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	if err != nil {
		log.Fatalf("Failed to create Grok instance: %s", err)
	}

	pattern := l.patternMap[logtype]
	fmt.Println("pattern =", pattern)
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
