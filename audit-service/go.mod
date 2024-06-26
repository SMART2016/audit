module audit-service

go 1.22

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elastic/go-elasticsearch/v8 v8.13.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/segmentio/kafka-go v0.4.47
	github.com/sirupsen/logrus v1.9.3
	github.com/vjeantet/grok v1.0.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/elastic/elastic-transport-go/v8 v8.5.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
