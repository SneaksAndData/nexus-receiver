package models

import coremodels "github.com/SneaksAndData/nexus-core/pkg/checkpoint/models"

type CheckRunResponse struct {
	IsProcessed bool `json:"is_processed"`
}

// FromCheckpoint converts a checkpoint to a response of a CheckRun endpoint
func FromCheckpoint(checkpoint *coremodels.CheckpointedRequest) *CheckRunResponse {
	return &CheckRunResponse{
		IsProcessed: checkpoint.IsFinished(),
	}
}
