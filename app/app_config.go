package app

import (
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"time"
)

type ReceiverConfig struct {
	CqlStore                   request.AstraBundleConfig `mapstructure:"cql-store,omitempty"`
	FailureRateBaseDelay       time.Duration             `mapstructure:"failure-rate-base-delay,omitempty"`
	FailureRateMaxDelay        time.Duration             `mapstructure:"failure-rate-max-delay,omitempty"`
	RateLimitElementsPerSecond int                       `mapstructure:"rate-limit-elements-per-second,omitempty"`
	RateLimitElementsBurst     int                       `mapstructure:"rate-limit-elements-burst,omitempty"`
	Workers                    int                       `mapstructure:"workers,omitempty"`
}
