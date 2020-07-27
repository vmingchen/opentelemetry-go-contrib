module go.opentelemetry.io/contrib/sdk/dynamicconfig

go 1.14

require (
	github.com/benbjohnson/clock v1.0.3
	github.com/open-telemetry/opentelemetry-collector v0.3.0
	github.com/open-telemetry/opentelemetry-proto v0.3.0
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.7.0
	google.golang.org/grpc v1.30.0
)

replace github.com/open-telemetry/opentelemetry-proto => github.com/vmingchen/opentelemetry-proto v0.3.1-0.20200716191220-7eb25882f08b

replace go.opentelemetry.io/contrib => ../../

replace go.opentelemetry.io/contrib/sdk/dynamicconfig => ./
