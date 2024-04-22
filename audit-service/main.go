package main

import (
	"net/http"

	auditlog "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
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
	logNormalizer := &LogNormalizer{make(map[string]string)}
	logPatternRegistration(logNormalizer)
	// Starting Audit-service API
	go http.ListenAndServe(":8080", Router{}.getRoutes())

	//Starting Audit service kafka consumer for log events from various sources.
	eventStore := getNewElasticsearchClient()
	storeLogEvents(logNormalizer, eventStore)
}
