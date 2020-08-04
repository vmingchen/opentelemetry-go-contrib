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

package basic

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/experimental/metricconfigservice"
	resourcepb "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	"go.opentelemetry.io/contrib/sdk/dynamicconfig/metric/controller/notify"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// A ServiceReader periodically reads from a remote configuration service to get configs that apply
// to the SDK.
type ServiceReader struct {
	clock clock.Clock // for testing

	conn   *grpc.ClientConn
	client pb.MetricConfigClient

	lastTimestamp        time.Time
	lastKnownFingerprint []byte
	resource             *resourcepb.Resource
}

func NewServiceReader(configHost string, resource *resourcepb.Resource) (*ServiceReader, error) {
	conn, err := grpc.Dial(configHost, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("fail to connect to config backend: %w", err)
	}

	return &ServiceReader{
		clock:    clock.New(),
		conn:     conn,
		client:   pb.NewMetricConfigClient(conn),
		resource: resource,
	}, nil
}

// ReadConfig reads the latest configuration data from the backend. Returns
// a nil *MetricConfig if there have been no changes to the configuration
// since the last check.
func (r *ServiceReader) ReadConfig() (*notify.MetricConfig, error) {
	request := &pb.MetricConfigRequest{
		LastKnownFingerprint: r.lastKnownFingerprint,
		Resource:             r.resource,
	}

	md := metadata.Pairs("timestamp", r.clock.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	response, err := r.client.GetMetricConfig(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("fail to get metric config: %w", err)
	}

	// TODO: SuggestedWaitTimeSec may not be read unless there is a change
	// reflected in the fingerprints
	if bytes.Equal(r.lastKnownFingerprint, response.Fingerprint) {
		return nil, nil
	}

	r.lastKnownFingerprint = response.Fingerprint
	r.lastTimestamp = r.clock.Now()

	newConfig := notify.MetricConfig{*response}
	if err := newConfig.Validate(); err != nil {
		return nil, fmt.Errorf("metric config invalid: %w", err)
	}

	return &newConfig, nil
}

func (r *ServiceReader) Stop() error {
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("fail to close connection to config backend: %w", err)
	}

	return nil
}