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
	"sync"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/contrib/exporters/metric/dynamicconfig"
	"go.opentelemetry.io/contrib/internal/transform"
)

// Wrapper used to verify we wait suggestedWaitTimeSec until we read
// from ServiceReader again.
type serviceReaderWrapper struct {
	reader *dynamicconfig.ServiceReader
	// Mutex for testVar.
	testLock sync.Mutex
	// testVar is used to test how long ServiceReader.ReadConfig()
	// is running.
	testVar int
}

func (r *serviceReaderWrapper) asyncReadConfig(t *testing.T) {
	r.testLock.Lock()
	r.testVar = 1
	r.testLock.Unlock()

	_, err := r.reader.ReadConfig()
	assert.NoError(t, err)

	r.testLock.Lock()
	r.testVar = 0
	r.testLock.Unlock()
}

func (r *serviceReaderWrapper) getTestVar() int {
	r.testLock.Lock()
	defer r.testLock.Unlock()
	return r.testVar
}

func TestReadConfig(t *testing.T) {
	clock := clock.NewMock()

	// Mock config service returns config with a suggested wait time of 5 minutes.
	config := dynamicconfig.GetDefaultConfig(60, DefaultFingerprint)
	config.SuggestedWaitTimeSec = 300
	stopFunc := runMockConfigService(t, ConfigServiceAddress, config)

	reader := dynamicconfig.NewServiceReader(
		ConfigServiceAddress,
		transform.Resource(mockResource("servicereadertest")),
	)
	reader.SetClock(clock)
	wrapper := serviceReaderWrapper{
		reader: reader,
	}

	// Sets ServiceReader.suggestedWaitTimeSec to 300 (5 minutes).
	_, err := reader.ReadConfig()
	assert.NoError(t, err)

	// Test that ServiceReader.ReadConfig() waits 5 minutes to read.
	go wrapper.asyncReadConfig(t)
	clock.Add(time.Minute)
	require.Equal(t, wrapper.getTestVar(), 1)
	clock.Add(4 * time.Minute)
	require.Equal(t, wrapper.getTestVar(), 0)

	stopFunc()
}
