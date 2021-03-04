module go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo

go 1.14

replace go.opentelemetry.io/contrib => ../../../..

require (
	github.com/labstack/echo/v4 v4.1.16
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib v0.10.0
	go.opentelemetry.io/otel v0.18.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout v0.10.0
	go.opentelemetry.io/otel/sdk v0.18.0
	go.opentelemetry.io/otel/sdk/metric v0.18.0 // indirect
	google.golang.org/grpc v1.31.0
)
