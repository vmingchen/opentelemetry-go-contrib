module go.opentelemetry.io/contrib/exporters/metric/dynamicconfig

go 1.14

replace go.opentelemetry.io/contrib => ../../..

replace go.opentelemetry.io/otel => github.com/open-telemetry/opentelemetry-go v0.6.1-0.20200623190015-2966505271c3

replace github.com/open-telemetry/opentelemetry-proto => ../../../../opentelemetry-proto

require (
	github.com/benbjohnson/clock v1.0.3
	github.com/open-telemetry/opentelemetry-collector v0.3.0
	github.com/open-telemetry/opentelemetry-proto v0.3.0
	github.com/stretchr/testify v1.4.0
	go.opentelemetry.io/contrib v0.6.1
	go.opentelemetry.io/otel v0.7.0
	google.golang.org/grpc v1.30.0
)