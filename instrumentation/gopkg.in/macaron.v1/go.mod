module go.opentelemetry.io/contrib/instrumentation/macaron

go 1.14

replace go.opentelemetry.io/contrib => ../../..

require (
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib v0.10.0
	go.opentelemetry.io/otel v0.17.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout v0.10.0
	go.opentelemetry.io/otel/sdk v0.17.0
	go.opentelemetry.io/otel/sdk/metric v0.17.0 // indirect
	gopkg.in/macaron.v1 v1.3.9
)
