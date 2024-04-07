package main

import (
	"encoding/json"
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

	g, err := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	if err != nil {
		log.Fatalf("Failed to create Grok instance: %s", err)
	}

	pattern := l.patternMap[logtype]
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
