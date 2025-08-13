package app

import (
	"fmt"
	coremodels "github.com/SneaksAndData/nexus-core/pkg/checkpoint/models"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"github.com/SneaksAndData/nexus-core/pkg/pipeline"
	"github.com/SneaksAndData/nexus-receiver/api/v1/models"
)

type CompletionActor = pipeline.DefaultPipelineStageActor[*models.CompletionInput, string]

func NewCompletionActor(store *request.CqlStore, appConfig *ReceiverConfig) *CompletionActor {
	return pipeline.NewDefaultPipelineStageActor[*models.CompletionInput, string](
		"request_completion",
		map[string]string{},
		appConfig.FailureRateBaseDelay,
		appConfig.FailureRateMaxDelay,
		appConfig.RateLimitElementsPerSecond,
		appConfig.RateLimitElementsBurst,
		appConfig.Workers,
		func(element *models.CompletionInput) (string, error) {
			return completeRequest(element, store)
		},
		nil,
	)
}

func completeRequest(input *models.CompletionInput, cqlStore *request.CqlStore) (string, error) {
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

		// TODO: metrics report
	} else {
		// TODO: metrics report
		requestCopy.LifecycleStage = coremodels.LifecycleStageFailed
		requestCopy.AlgorithmFailureCause = input.Result.ErrorCause
		requestCopy.AlgorithmFailureDetails = input.Result.ErrorDetails
	}

	insertErr := cqlStore.UpsertCheckpoint(requestCopy)

	if insertErr != nil { // coverage-ignore
		return "", insertErr
	}

	return requestCopy.Id, nil
}
