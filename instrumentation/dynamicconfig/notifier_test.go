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

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/contrib/instrumentation/dynamicconfig"

	controllerTest "go.opentelemetry.io/otel/sdk/metric/controller/test"
)

const CONFIG_SERVICE_ADDRESS = "localhost:50420"
var DEFAULT_FINGERPRINT = []byte{'f', 'o', 'o'}

// testLock is to prevent race conditions in test code
// testVar is used to verify OnInitialConfig and OnUpdatedConfig are called
type testWatcher struct {
	testLock sync.Mutex
	testVar  int
}

func (w *testWatcher) OnInitialConfig(config *dynamicconfig.Config) {
	w.testLock.Lock()
	defer w.testLock.Unlock()
	w.testVar = 1
}

func (w *testWatcher) OnUpdatedConfig(config *dynamicconfig.Config) {
	w.testLock.Lock()
	defer w.testLock.Unlock()
	w.testVar = 2
}

// Use a getter to prevent race conditions around testVar
func (w *testWatcher) getTestVar() int {
	w.testLock.Lock()
	defer w.testLock.Unlock()
	return w.testVar
}

func newExampleNotifier() *dynamicconfig.Notifier {
	notifier, _ := dynamicconfig.NewNotifier(
		30 * time.Second,
		dynamicconfig.GetDefaultConfig(30, []byte{'b', 'a', 'r'}),
		dynamicconfig.WithConfigHost(CONFIG_SERVICE_ADDRESS),
		dynamicconfig.WithResource(mockResource("notifiertest")),
	)

	return notifier
}

// Test config updates
func TestDynamicNotifier(t *testing.T) {
	watcher := testWatcher{
		testVar: 0,
	}
	mock := controllerTest.NewMockClock()

	stopFunc := runMockConfigService(
		t,
		CONFIG_SERVICE_ADDRESS,
		dynamicconfig.GetDefaultConfig(60, DEFAULT_FINGERPRINT),
	)

	notifier := newExampleNotifier()
	require.Equal(t, watcher.getTestVar(), 0)

	notifier.SetClock(mock)
	notifier.Start()

	notifier.Register(&watcher)
	require.Equal(t, watcher.getTestVar(), 1)

	mock.Add(time.Minute)

	require.Equal(t, watcher.getTestVar(), 2)

	notifier.Stop()
	stopFunc()
}

// Test config doesn't update
func TestNonDynamicNotifier(t *testing.T) {
	watcher := testWatcher{
		testVar: 0,
	}
	mock := controllerTest.NewMockClock()
	notifier, _ := dynamicconfig.NewNotifier(
		time.Minute,
		dynamicconfig.GetDefaultConfig(60, DEFAULT_FINGERPRINT),
	)
	require.Equal(t, watcher.getTestVar(), 0)

	notifier.SetClock(mock)
	notifier.Start()

	notifier.Register(&watcher)
	require.Equal(t, watcher.getTestVar(), 1)

	mock.Add(time.Minute)

	require.Equal(t, watcher.getTestVar(), 1)
	notifier.Stop()
}

func TestDoubleStop(t *testing.T) {
	stopFunc := runMockConfigService(
		t,
		CONFIG_SERVICE_ADDRESS,
		dynamicconfig.GetDefaultConfig(60, DEFAULT_FINGERPRINT),
	)
	notifier := newExampleNotifier()
	notifier.Start()
	notifier.Stop()
	notifier.Stop()
	stopFunc()
}

func TestPushDoubleStart(t *testing.T) {
	stopFunc := runMockConfigService(
		t,
		CONFIG_SERVICE_ADDRESS,
		dynamicconfig.GetDefaultConfig(60, DEFAULT_FINGERPRINT),
	)
	notifier := newExampleNotifier()
	notifier.Start()
	notifier.Start()
	notifier.Stop()
	stopFunc()
}
