module go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama

go 1.14

replace go.opentelemetry.io/contrib => ../../../..

require (
	github.com/Shopify/sarama v1.26.4
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib v0.10.0
	go.opentelemetry.io/otel v0.17.0 // indirect
	google.golang.org/grpc v1.31.0
)
