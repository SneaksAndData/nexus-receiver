package models

// CompletionInput represents the data processed by completion actor
type CompletionInput struct {
	Result        AlgorithmResult
	RequestId     string
	AlgorithmName string
}
