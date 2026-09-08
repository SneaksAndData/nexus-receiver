package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/SneaksAndData/nexus-receiver/api/v1/models"
)

const baseURL = "http://localhost:5555/receiver"

func TestSmoke_CheckRun_InitialState(t *testing.T) {
	algorithmName := "test-algorithm"
	requestId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	url := fmt.Sprintf("%s/algorithm/v1/check/%s/requests/%s", baseURL, algorithmName, requestId)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to perform GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200 OK, got %d. body: %s", resp.StatusCode, string(body))
	}

	var checkResp models.CheckRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		t.Fatalf("failed to decode check run response: %v", err)
	}

	// In initial seed data, lifecycle_stage is RUNNING, so is_processed should be false
	if checkResp.IsProcessed {
		t.Errorf("expected is_processed to be false for initial RUNNING checkpoint, got true")
	}
}

func TestSmoke_CheckRun_NotFound(t *testing.T) {
	algorithmName := "non-existent-algorithm"
	requestId := "00000000-0000-0000-0000-000000000000"

	url := fmt.Sprintf("%s/algorithm/v1/check/%s/requests/%s", baseURL, algorithmName, requestId)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to perform GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 404 Not Found, got %d. body: %s", resp.StatusCode, string(body))
	}
}

func TestSmoke_CompleteRun_And_CheckRun(t *testing.T) {
	algorithmName := "test-algorithm"
	requestId := "2c7b6e8d-cc3c-4b5b-a3f6-5d7b9e2c7f2a"

	payload := models.AlgorithmResult{
		ResultUri: "http://localhost:9000/smoke-test-result",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal complete payload: %v", err)
	}

	completeURL := fmt.Sprintf("%s/algorithm/v1/complete/%s/requests/%s", baseURL, algorithmName, requestId)
	resp, err := http.Post(completeURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("failed to perform POST %s: %v", completeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 202 Accepted, got %d. body: %s", resp.StatusCode, string(body))
	}

	// Completion actor processes asynchronously in the background.
	// Poll check endpoint until is_processed becomes true or timeout occurs.
	checkURL := fmt.Sprintf("%s/algorithm/v1/check/%s/requests/%s", baseURL, algorithmName, requestId)
	deadline := time.Now().Add(15 * time.Second)
	completed := false

	for time.Now().Before(deadline) {
		time.Sleep(500 * time.Millisecond)

		checkRespHttp, err := http.Get(checkURL)
		if err != nil {
			t.Logf("check request error: %v, retrying...", err)
			continue
		}

		if checkRespHttp.StatusCode != http.StatusOK {
			checkRespHttp.Body.Close()
			continue
		}

		var checkResult models.CheckRunResponse
		decodeErr := json.NewDecoder(checkRespHttp.Body).Decode(&checkResult)
		checkRespHttp.Body.Close()

		if decodeErr == nil && checkResult.IsProcessed {
			completed = true
			break
		}
	}

	if !completed {
		t.Fatalf("timed out waiting for checkpoint %s to be marked as processed", requestId)
	}
}

func TestSmoke_CompleteRun_InvalidPayload(t *testing.T) {
	algorithmName := "test-algorithm"
	requestId := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	completeURL := fmt.Sprintf("%s/algorithm/v1/complete/%s/requests/%s", baseURL, algorithmName, requestId)
	resp, err := http.Post(completeURL, "application/json", strings.NewReader("invalid-json"))
	if err != nil {
		t.Fatalf("failed to perform POST %s: %v", completeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 400 Bad Request for malformed payload, got %d. body: %s", resp.StatusCode, string(body))
	}
}
