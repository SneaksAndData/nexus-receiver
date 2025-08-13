package app

import (
	"context"
	coremodels "github.com/SneaksAndData/nexus-core/pkg/checkpoint/models"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"github.com/SneaksAndData/nexus-core/pkg/pipeline"
	"github.com/SneaksAndData/nexus-receiver/api/v1/models"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/ktesting"
	"testing"
	"time"
)

type fixture struct {
	actor    *CompletionActor
	cqlStore *request.CqlStore
	ctx      context.Context
}

func newFixture(t *testing.T) *fixture {
	_, ctx := ktesting.NewTestContext(t)
	f := &fixture{}

	f.ctx = ctx
	f.cqlStore = request.NewScyllaCqlStore(
		klog.FromContext(ctx), &request.ScyllaCqlStoreConfig{
			Hosts: []string{"127.0.0.1"},
		})
	f.actor = NewCompletionActor(f.ctx, f.cqlStore, &ReceiverConfig{
		AstraCqlStore:              request.AstraBundleConfig{},
		ScyllaCqlStore:             request.ScyllaCqlStoreConfig{},
		CqlStoreType:               "scylla",
		FailureRateBaseDelay:       time.Second,
		FailureRateMaxDelay:        time.Second * 2,
		RateLimitElementsPerSecond: 10,
		RateLimitElementsBurst:     10,
		Workers:                    2,
		LogLevel:                   "INFO",
		BindPort:                   0,
	})

	return f
}

func TestCompletion(t *testing.T) {
	f := newFixture(t)

	go f.actor.Start(f.ctx, pipeline.NewActorPostStart(func(ctx context.Context) error {
		f.actor.Receive(&models.CompletionInput{
			Result: models.AlgorithmResult{
				ResultUri: "http://localhost:9000/data",
			},
			RequestId:     "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			AlgorithmName: "test-algorithm",
		})
		f.actor.Receive(&models.CompletionInput{
			Result: models.AlgorithmResult{
				ErrorCause:   "Failed run",
				ErrorDetails: "Unit test generated a failure here.",
			},
			RequestId:     "2c7b6e8d-cc3c-4b5b-a3f6-5d7b9e2c7f2a",
			AlgorithmName: "test-algorithm",
		})
		return nil
	}))

	// allow pipeline to initialize and run
	time.Sleep(time.Second * 5)

	results, err := f.cqlStore.ReadCheckpointsByTag("tag_123")

	if err != nil {
		t.Errorf("error while reading checkpoints: %v", err)
		t.FailNow()
	}

	for result := range results {
		if result.Id == "f47ac10b-58cc-4372-a567-0e02b2c3d479" && (result.ResultUri != "http://localhost:9000/data" || result.LifecycleStage != coremodels.LifecycleStageCompleted) {
			t.Errorf("expected a completed checkpoint with result url http://localhost:9000/data, but got result url: %s, lifecycle stage: %s", result.ResultUri, result.LifecycleStage)
			t.FailNow()
		}
		if result.Id == "2c7b6e8d-cc3c-4b5b-a3f6-5d7b9e2c7f2a" && (result.ResultUri != "" || result.LifecycleStage != coremodels.LifecycleStageFailed) {
			t.Errorf("expected a failed checkpoint without a result url, but got result url: %s, lifecycle stage: %s", result.ResultUri, result.LifecycleStage)
			t.FailNow()
		}
	}
}
