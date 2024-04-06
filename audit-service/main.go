package main

import (
	auditlog "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
)

const (
	USER_SERVICE_LOG_TYPE    = "user-service"
	USER_SERVICE_LOG_PATTERN = `CurrentUser: %{WORD:currentUser}, System: %{USERNAME:system}, Action: %{WORD:action}, IP: \[%{IPV6:ip}\]:%{NUMBER:port}, Agent: %{GREEDYDATA:agent}, Time: %{TIMESTAMP_ISO8601:time}`
)

func init() {
	//file, err := os.OpenFile("service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	auditlog.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/service.log", // Log file path
		MaxSize:    1,                    // Maximum size of a log file before rotation (in megabytes)
		MaxBackups: 3,                    // Maximum number of old log files to retain
		MaxAge:     28,                   // Maximum number of days to retain an old log file
		Compress:   true,                 // Compress/zip old log files
	})

	//auditlog.SetOutput(file)
	// Set log level
	auditlog.SetLevel(auditlog.InfoLevel)

	// Use JSON formatter
	auditlog.SetFormatter(&auditlog.JSONFormatter{})
}

func main() {
	logNormalizer := LogNormalizer{make(map[string]string)}
	logNormalizer.registerLogPatterns(USER_SERVICE_LOG_TYPE, USER_SERVICE_LOG_PATTERN)
	// Starting Audit-service API
	go http.ListenAndServe(":9191", Router{}.getRoutes())

	//Starting Audit service kafka consumer for log events from various sources.
	brokers := []string{"localhost:9093"}
	topic := "log_events_topic"
	logsChan := consumeKafkaMessages(brokers, topic)
	for logMsg := range logsChan {
		// Normalize your log message
		normalizedLog := logNormalizer.normalizeLog(USER_SERVICE_LOG_TYPE, logMsg)
		auditlog.Info(normalizedLog)
		getNewElasticsearchClient().pushLogEvents(normalizedLog)
	}

}
