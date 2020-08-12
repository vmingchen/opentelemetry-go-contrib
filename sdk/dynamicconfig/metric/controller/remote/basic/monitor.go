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
	"time"

	pb "go.opentelemetry.io/contrib/sdk/dynamicconfig/internal/proto/experimental/metrics/configservice"
	"go.opentelemetry.io/contrib/sdk/dynamicconfig/internal/transform"
	"go.opentelemetry.io/contrib/sdk/dynamicconfig/metric/controller/remote"
	controllerTime "go.opentelemetry.io/otel/sdk/metric/controller/time"
	"go.opentelemetry.io/otel/sdk/resource"
)

const DefaultCheckFrequency = 30 * time.Minute

// A Monitor monitors a config service for a config changing, then letting
// all its subscribers know if the config has changed.
//
// All fields except for subscribed and config, which are protected by lock,
// should be read-only once set.
type Monitor struct {
	clock        controllerTime.Clock
	configHost   string
	lastWaitTime int32
	resource     *resource.Resource
	ticker       controllerTime.Ticker
}

// Constructor for a Monitor
func NewMonitor(configHost string, resource *resource.Resource) *Monitor {
	monitor := &Monitor{
		clock:      controllerTime.RealClock{},
		configHost: configHost,
		resource:   resource,
	}

	return monitor
}

func (m *Monitor) MonitorChanges(mch remote.MonitorChannel) {
	go func() {
		m.ticker = m.clock.Ticker(DefaultCheckFrequency)
		serviceReader, err := NewServiceReader(m.configHost, transform.Resource(m.resource))
		if err != nil {
			mch.Err <- err
		}

		m.tick(mch.Data, mch.Err, serviceReader)
		for {
			select {
			case <-m.ticker.C():
				m.tick(mch.Data, mch.Err, serviceReader)

			case <-mch.Quit:
				m.ticker.Stop()
				if err := serviceReader.Stop(); err != nil {
					mch.Err <- err
				}
				return
			}
		}
	}()
}

func (m *Monitor) tick(data chan<- []*pb.MetricConfigResponse_Schedule, errCh chan<- error, serviceReader *ServiceReader) {
	newConfig, err := serviceReader.ReadConfig()
	if err != nil {
		errCh <- err
	}

	if newConfig != nil {
		m.updateWaitTime(newConfig.SuggestedWaitTimeSec)
		data <- newConfig.Schedules
	}
}

func (m *Monitor) updateWaitTime(waitTime int32) {
	if waitTime > 0 && m.lastWaitTime != waitTime {
		m.ticker.Stop()
		m.lastWaitTime = waitTime
		m.ticker = m.clock.Ticker(time.Duration(m.lastWaitTime) * time.Second)
	}
}
