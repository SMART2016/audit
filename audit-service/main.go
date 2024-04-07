package main

import (
	auditlog "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
)

const (
	USER_SERVICE_LOG_TYPE  = "user-service"
	AUTH_SERVICE_LOG_TYPE  = "auth-service"
	AUDIT_SERVICE_LOG_TYPE = "audit-service"
	SERVICE_LOG_PATTERN    = `RequestId: %{DATA:RequestId}, CurrentUser: %{DATA:CurrentUser},Role: %{DATA:Role}, System: %{DATA:System}, Action: %{DATA:Action}, IP: \[%{IP:IP}\]:%{NUMBER:Port}, Agent: %{DATA:Agent}, Time: %{TIMESTAMP_ISO8601:Time}, Status: %{DATA:Status}
`
)

var logger lumberjack.Logger

func init() {
	//Log rotation with lumberjack
	logger = lumberjack.Logger{
		Filename:   "./logs/service.log", // Log file path
		MaxSize:    1,                    // Maximum size of a log file before rotation (in megabytes)
		MaxBackups: 3,                    // Maximum number of old log files to retain
		MaxAge:     28,                   // Maximum number of days to retain an old log file
		Compress:   true,                 // Compress/zip old log files
	}
	auditlog.SetOutput(&logger)

	// Set log level
	auditlog.SetLevel(auditlog.InfoLevel)

	// Use JSON formatter
	auditlog.SetFormatter(&auditlog.JSONFormatter{})

}
func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the main handler\n"))
}
func main() {
	// register Pattern for all Log types
	logNormalizer := LogNormalizer{make(map[string]string)}
	logPatternRegistration(logNormalizer)
	// Starting Audit-service API
	go http.ListenAndServe(":9191", Router{}.getRoutes())

	//Starting Audit service kafka consumer for log events from various sources.
	startKafkaConsumer(logNormalizer)
}

func logPatternRegistration(normalizer LogNormalizer) {
	normalizer.registerLogPatterns(USER_SERVICE_LOG_TYPE, SERVICE_LOG_PATTERN)
	normalizer.registerLogPatterns(AUTH_SERVICE_LOG_TYPE, SERVICE_LOG_PATTERN)
	normalizer.registerLogPatterns(AUDIT_SERVICE_LOG_TYPE, SERVICE_LOG_PATTERN)
}
