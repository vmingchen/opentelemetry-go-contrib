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

package push

import (
	pb "go.opentelemetry.io/contrib/sdk/dynamicconfig/internal/proto/experimental/metrics/configservice"
	"go.opentelemetry.io/contrib/sdk/dynamicconfig/metric/controller/remote"
	"go.opentelemetry.io/contrib/sdk/dynamicconfig/metric/controller/remote/mock"
	controllerTime "go.opentelemetry.io/otel/sdk/metric/controller/time"
)

func (c *Controller) SetClock(clock controllerTime.Clock) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.clock = clock
}

func (c *Controller) SetMonitor(monitor remote.Monitor) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.monitor = monitor
}

func (c *Controller) SetPeriod(period int32) {
	scheds := []*pb.MetricConfigResponse_Schedule{
		{
			InclusionPatterns: []*pb.MetricConfigResponse_Schedule_Pattern{
				{
					Match: &pb.MetricConfigResponse_Schedule_Pattern_StartsWith{
						StartsWith: "",
					},
				},
			},
			PeriodSec: period,
		},
	}

	monitor := mock.NewMonitor()
	monitor.Receive(scheds)
	c.SetMonitor(monitor)
}

func (c *Controller) SetDone() {
	c.done = make(chan struct{})
}

func (c *Controller) WaitDone() {
	<-c.done
}
