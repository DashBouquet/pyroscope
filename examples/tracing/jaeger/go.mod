module rideshare

go 1.17

require (
	github.com/pyroscope-io/client v0.2.0
	github.com/pyroscope-io/otelpyroscope v0.1.0
	github.com/sirupsen/logrus v1.8.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/trace v1.4.1
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
)
