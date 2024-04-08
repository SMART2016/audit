package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vjeantet/grok"
)

const (
	USER_SERVICE_LOG_TYPE  = "user-service"
	AUTH_SERVICE_LOG_TYPE  = "auth-service"
	AUDIT_SERVICE_LOG_TYPE = "audit-service"
	SERVICE_LOG_PATTERN    = `RequestId: %{DATA:RequestId}, CurrentUser: %{DATA:CurrentUser},Role: %{DATA:Role}, System: %{DATA:System}, Action: %{DATA:Action}, IP: (?:\[%{IP:IP}\]|%{IP:IP}):%{NUMBER:Port}, Agent: %{DATA:Agent}, Time: %{TIMESTAMP_ISO8601:Time}, Status: %{DATA:Status}
`

	//			`RequestId: %{DATA:RequestId}, CurrentUser: %{DATA:CurrentUser},Role: %{DATA:Role}, System: %{DATA:System}, Action: %{DATA:Action}, IP: \[%{IP:IP}\]:%{NUMBER:Port}, Agent: %{DATA:Agent}, Time: %{TIMESTAMP_ISO8601:Time}, Status: %{DATA:Status}
	//`
)

type LogNormalizer struct {
	patternMap map[string]string
}

func (l *LogNormalizer) RegisterLogPatterns(logtype string, pattern string) {
	l.patternMap[logtype] = pattern
}
func (l *LogNormalizer) NormalizeLog(logtype string, logMsg string) string {
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

func logPatternRegistration(normalizer EventLogNormalizer) {
	normalizer.RegisterLogPatterns(USER_SERVICE_LOG_TYPE, SERVICE_LOG_PATTERN)
	normalizer.RegisterLogPatterns(AUTH_SERVICE_LOG_TYPE, SERVICE_LOG_PATTERN)
	normalizer.RegisterLogPatterns(AUDIT_SERVICE_LOG_TYPE, SERVICE_LOG_PATTERN)
}
