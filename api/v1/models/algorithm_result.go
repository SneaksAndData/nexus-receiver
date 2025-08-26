package models

// AlgorithmResult contains an optional URI to download results, and error information in case of a failed run.
type AlgorithmResult struct {
	// URL to download results.
	ResultUri string `json:"resultUri"`

	// Failure cause, if any.
	ErrorCause string `json:"errorCause,omitempty"`

	// Failure details, if any.
	ErrorDetails string `json:"errorDetails,omitempty"`
}
