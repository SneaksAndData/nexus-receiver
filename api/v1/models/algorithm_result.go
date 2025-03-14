package models

type AlgorithmResult struct {
	// URL to download results.
	SasUri string `json:"sasUri"`

	// Failure cause, if any.
	Cause string `json:"cause,omitempty"`

	// Failure details, if any.
	Message string `json:"message,omitempty"`

	// Failure error code, if any.
	ErrorCode string `json:"errorCode,omitempty"`
}
