package models

// CompletionInput represents the payload sent to the run completion endpoint.
type CompletionInput struct {
	Result        AlgorithmResult
	RequestId     string
	AlgorithmName string
}
