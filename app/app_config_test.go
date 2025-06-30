package app

import (
	"context"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	nexusconf "github.com/SneaksAndData/nexus-core/pkg/configurations"
	"os"
	"reflect"
	"testing"
	"time"
)

func getExpectedConfig() *ReceiverConfig {
	return &ReceiverConfig{
		CqlStore: request.AstraBundleConfig{
			SecureConnectionBundleBase64: "base64value",
			GatewayUser:                  "user",
			GatewayPassword:              "password",
		},
		FailureRateBaseDelay:       time.Millisecond * 100,
		FailureRateMaxDelay:        time.Second,
		RateLimitElementsPerSecond: 10,
		RateLimitElementsBurst:     100,
		Workers:                    10,
		LogLevel:                   "INFO",
		BindPort:                   8080,
	}
}

func Test_LoadConfig(t *testing.T) {
	var expected = getExpectedConfig()

	var result = nexusconf.LoadConfig[ReceiverConfig](context.TODO())
	if !reflect.DeepEqual(*expected, result) {
		t.Errorf("LoadConfig failed, expected %v, got %v", *expected, result)
	}
}

func Test_LoadConfigFromEnv(t *testing.T) {
	_ = os.Setenv("NEXUS__CQL_STORE__GATEWAY_PASSWORD", "password1")
	var expected = getExpectedConfig()
	expected.CqlStore.GatewayPassword = "password1"

	var result = nexusconf.LoadConfig[ReceiverConfig](context.TODO())
	if !reflect.DeepEqual(*expected, result) {
		t.Errorf("LoadConfig failed, expected %v, got %v", *expected, result)
	}
}
