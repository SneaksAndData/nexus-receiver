package app

import (
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/store/cassandra"

	"time"
)

type ReceiverConfig struct {
	AstraCqlStore              cassandra.AstraBundleConfig `mapstructure:"astra-cql-store,omitempty"`
	ScyllaCqlStore             cassandra.ScyllaConfig      `mapstructure:"scylla-cql-store,omitempty"`
	KeyspacesCqlStore          cassandra.KeyspacesConfig   `mapstructure:"keyspaces-cql-store,omitempty"`
	CqlStoreType               string                      `mapstructure:"cql-store-type,omitempty"`
	FailureRateBaseDelay       time.Duration               `mapstructure:"failure-rate-base-delay,omitempty"`
	FailureRateMaxDelay        time.Duration               `mapstructure:"failure-rate-max-delay,omitempty"`
	RateLimitElementsPerSecond int                         `mapstructure:"rate-limit-elements-per-second,omitempty"`
	RateLimitElementsBurst     int                         `mapstructure:"rate-limit-elements-burst,omitempty"`
	Workers                    int                         `mapstructure:"workers,omitempty"`
	LogLevel                   string                      `mapstructure:"log-level,omitempty"`
	BindPort                   int                         `mapstructure:"bind-port,omitempty"`
}

const (
	CqlStoreAstra     = "cassandra-astra"
	CqlStoreScylla    = "cassandra-scylla"
	CqlStoreKeyspaces = "cassandra-keyspaces"
)
