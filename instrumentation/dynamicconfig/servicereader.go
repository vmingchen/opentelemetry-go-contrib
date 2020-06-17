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

package dynamicconfig

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	resourcepb "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	pb "github.com/vmingchen/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

type ServiceReader struct {
	// Required
	configHost string

	// Timestamp of last time config service was checked
	lastTimestamp time.Time

	// Most recent config version
	lastKnownFingerprint []byte

	// Suggested time to wait before checking config service again (seconds)
	suggestedWaitTimeSec int32

	// Required. Label to identify this instance
	resource *resourcepb.Resource
}

func newServiceReader(configHost string, resource *resourcepb.Resource) *ServiceReader {
	return &ServiceReader{
		configHost:   configHost,
		resource:     resource,
	}
}

// Reads from a config service
func (r *ServiceReader) readConfig() (*Config, error) {
	// Wait for the suggestedWaitTimeSec
	if !r.lastTimestamp.IsZero() && r.suggestedWaitTimeSec != 0 {
		suggestedReadTime := r.lastTimestamp.Add(time.Duration(r.suggestedWaitTimeSec) * time.Second)
		time.Sleep(suggestedReadTime.Sub(time.Now()))

		r.suggestedWaitTimeSec = 0
	}

	// Get the new config
	conn, err := grpc.Dial(r.configHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb.NewDynamicConfigClient(conn)

	request := &pb.ConfigRequest{
		LastKnownFingerprint: r.lastKnownFingerprint,
		Resource:        r.resource,
	}

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	response, err := c.GetConfig(ctx, request)
	if err != nil {
		return nil, err
	}

	r.lastKnownFingerprint = response.Fingerprint
	r.lastTimestamp = time.Now()
	r.suggestedWaitTimeSec = response.SuggestedWaitTimeSec

	newConfig := Config{
		Fingerprint: response.Fingerprint,
		MetricConfig: response.MetricConfig,
		TraceConfig:  response.TraceConfig,
	}

	return &newConfig, nil
}
