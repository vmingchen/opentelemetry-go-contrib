// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynamicconfig_test

import (
	"context"
	"net"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector/translator/conventions"
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/exporters/metric/dynamicconfig"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

type mockServer struct {
	pb.UnimplementedDynamicConfigServer
	config *dynamicconfig.Config
}

// GetConfig implemented DynamicConfigServer
func (s *mockServer) GetConfig(ctx context.Context, in *pb.ConfigRequest) (*pb.ConfigResponse, error) {
	config := &pb.ConfigResponse{
		MetricConfig: s.config.MetricConfig,
		TraceConfig: s.config.TraceConfig,
	}
	return config, nil
}

// This function runs a mock config service at an address, serving a defined config.
// It returns a callback that stops the service.
func runMockConfigService(t *testing.T, addr string, config *dynamicconfig.Config) func() {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to get an address: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterDynamicConfigServer(srv, &mockServer{ config: config })

	go func() {
		_ = srv.Serve(ln)
	}()

	return func() {
		srv.Stop()
		_ = ln.Close()
	}
}

func mockResource(serviceName string) *resource.Resource {
	return resource.New(kv.Key(conventions.AttributeServiceName).String(serviceName))
}