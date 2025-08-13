package app

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-go/v5/statsd"
	coremodels "github.com/SneaksAndData/nexus-core/pkg/checkpoint/models"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"github.com/SneaksAndData/nexus-core/pkg/pipeline"
	"github.com/SneaksAndData/nexus-core/pkg/telemetry"
	"github.com/SneaksAndData/nexus-receiver/api/v1/models"
	"time"
)

type CompletionActor = pipeline.DefaultPipelineStageActor[*models.CompletionInput, string]

func NewCompletionActor(ctx context.Context, store *request.CqlStore, appConfig *ReceiverConfig) *CompletionActor {
	return pipeline.NewDefaultPipelineStageActor[*models.CompletionInput, string](
		"request_completion",
		map[string]string{},
		appConfig.FailureRateBaseDelay,
		appConfig.FailureRateMaxDelay,
		appConfig.RateLimitElementsPerSecond,
		appConfig.RateLimitElementsBurst,
		appConfig.Workers,
		func(element *models.CompletionInput) (string, error) {
			return completeRequest(element, store, telemetry.GetClient(ctx))
		},
		nil,
	)
}

func completeRequest(input *models.CompletionInput, cqlStore *request.CqlStore, metrics *statsd.Client) (string, error) {
	if input == nil { // coverage-ignore
		return "", fmt.Errorf("buffer is nil")
	}

	requestToComplete, err := cqlStore.ReadCheckpoint(input.AlgorithmName, input.RequestId)

	if err != nil { // coverage-ignore
		return "", err
	}

	if requestToComplete.IsFinished() {
		return requestToComplete.Id, nil
	}

	requestCopy := requestToComplete.DeepCopy()

	if input.Result.ErrorCause == "" {
		requestCopy.LifecycleStage = coremodels.LifecycleStageCompleted
		requestCopy.ResultUri = input.Result.ResultUri

		telemetry.Increment(metrics, "completions", map[string]string{"algorithm": requestCopy.Algorithm})
	} else {
		telemetry.Increment(metrics, "failures", map[string]string{"algorithm": requestCopy.Algorithm})
		requestCopy.LifecycleStage = coremodels.LifecycleStageFailed
		requestCopy.AlgorithmFailureCause = input.Result.ErrorCause
		requestCopy.AlgorithmFailureDetails = input.Result.ErrorDetails
	}

	// record run duration from receive to finish
	telemetry.Gauge(metrics, "run_duration", time.Since(requestCopy.ReceivedAt).Seconds(), map[string]string{"algorithm": requestCopy.Algorithm, "final_stage": requestCopy.LifecycleStage}, 1)

	insertErr := cqlStore.UpsertCheckpoint(requestCopy)

	if insertErr != nil { // coverage-ignore
		return "", insertErr
	}

	return requestCopy.Id, nil
}
