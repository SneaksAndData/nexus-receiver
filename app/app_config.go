package app

import (
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"time"
)

type ReceiverConfig struct {
	AstraCqlStore              request.AstraBundleConfig    `mapstructure:"astra-cql-store,omitempty"`
	ScyllaCqlStore             request.ScyllaCqlStoreConfig `mapstructure:"scylla-cql-store,omitempty"`
	CqlStoreType               string                       `mapstructure:"cql-store-type,omitempty"`
	FailureRateBaseDelay       time.Duration                `mapstructure:"failure-rate-base-delay,omitempty"`
	FailureRateMaxDelay        time.Duration                `mapstructure:"failure-rate-max-delay,omitempty"`
	RateLimitElementsPerSecond int                          `mapstructure:"rate-limit-elements-per-second,omitempty"`
	RateLimitElementsBurst     int                          `mapstructure:"rate-limit-elements-burst,omitempty"`
	Workers                    int                          `mapstructure:"workers,omitempty"`
	LogLevel                   string                       `mapstructure:"log-level,omitempty"`
	BindPort                   int                          `mapstructure:"bind-port,omitempty"`
}

const (
	CqlStoreAstra  = "astra"
	CqlStoreScylla = "scylla"
)
