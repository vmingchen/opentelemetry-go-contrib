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

package notifier

import (
	"bytes"
	"errors"

	pb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/dynamicconfig/v1"
)

// A configuration used in the SDK to dynamically change metric collection and tracing.
type Config struct {
	pb.ConfigResponse
}

// This is for convenient development/testing purposes.
// It produces a Config with a schedule that matches all instruments, with a
// collection period of `period`
func GetDefaultConfig(period pb.ConfigResponse_MetricConfig_Schedule_CollectionPeriod, fingerprint []byte) *Config {
	pattern := pb.ConfigResponse_MetricConfig_Schedule_Pattern{
		Match: &pb.ConfigResponse_MetricConfig_Schedule_Pattern_StartsWith{StartsWith: "*"},
	}
	schedule := pb.ConfigResponse_MetricConfig_Schedule{
		InclusionPatterns: []*pb.ConfigResponse_MetricConfig_Schedule_Pattern{&pattern},
		Period:            period,
	}

	return &Config{
		pb.ConfigResponse{
			Fingerprint: fingerprint,
			MetricConfig: &pb.ConfigResponse_MetricConfig{
				Schedules: []*pb.ConfigResponse_MetricConfig_Schedule{&schedule},
			},
		},
	}
}

func (config *Config) ValidateMetricConfig() error {
	if config.MetricConfig == nil {
		return errors.New("No MetricConfig")
	}

	for _, schedule := range config.MetricConfig.Schedules {
		if schedule.Period < 0 {
			return errors.New("Periods must be positive")
		}
	}

	return nil
}

func (config *Config) Equals(otherConfig *Config) bool {
	return bytes.Equal(config.Fingerprint, otherConfig.Fingerprint)
}
