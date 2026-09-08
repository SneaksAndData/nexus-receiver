package app

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/store/cassandra"
	nexusconf "github.com/SneaksAndData/nexus-core/pkg/configurations"
)

func getExpectedConfig() *ReceiverConfig {
	return &ReceiverConfig{
		AstraCqlStore: cassandra.AstraBundleConfig{
			SecureConnectionBundleBase64: "base64value",
			GatewayUser:                  "user",
			GatewayPassword:              "password",
			IndexesSupported:             false,
			Keyspace:                     "nexus",
		},
		ScyllaCqlStore: cassandra.ScyllaConfig{
			Hosts:            []string{"host1", "host2"},
			IndexesSupported: true,
			Keyspace:         "nexus",
		},
		KeyspacesCqlStore: cassandra.KeyspacesConfig{
			Keyspace: "nexus",
			Hosts:    []string{"keyspaces.aws.com"},
			Port:     "9042",
			CaPath:   "/tmp/ca",
			Region:   "us-east-1",
			UseIRSA:  false,
		},
		CqlStoreType:               CqlStoreAstra,
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
	_ = os.Setenv("NEXUS__ASTRA_CQL_STORE__GATEWAY_PASSWORD", "password1")
	var expected = getExpectedConfig()
	expected.AstraCqlStore.GatewayPassword = "password1"

	var result = nexusconf.LoadConfig[ReceiverConfig](context.TODO())
	if !reflect.DeepEqual(*expected, result) {
		t.Errorf("LoadConfig failed, expected %v, got %v", *expected, result)
	}
}
