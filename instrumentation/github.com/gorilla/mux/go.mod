module go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux

go 1.14

replace go.opentelemetry.io/contrib => ../../../..

require (
	github.com/DataDog/sketches-go v0.0.1 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib v0.10.0
	go.opentelemetry.io/otel v0.17.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
)
