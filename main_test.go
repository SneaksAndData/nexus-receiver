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

	coremodels "github.com/SneaksAndData/nexus-core/pkg/checkpoint/models"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/store"
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/store/cassandra"
	"github.com/SneaksAndData/nexus-receiver/api/v1/models"
	"github.com/google/uuid"
	"k8s.io/klog/v2"
)

const (
	baseURL               = "http://localhost:5555/receiver"
	templateAlgorithmName = "test-algorithm"
	templateCheckpointID  = "f47ac10b-58cc-4372-a567-0e02b2c3d479"
)

func newTestStore(t *testing.T) store.CheckpointStore {
	return cassandra.NewScyllaStore(
		klog.FromContext(t.Context()),
		&cassandra.ScyllaConfig{
			Hosts:            []string{"127.0.0.1"},
			Port:             "30042",
			Keyspace:         "nexus",
			IndexesSupported: true,
		},
	)
}

func provisionCheckpoint(t *testing.T, cqlStore store.CheckpointStore) *coremodels.CheckpointedRequest {
	template, err := cqlStore.ReadCheckpoint(templateAlgorithmName, templateCheckpointID)
	if err != nil {
		t.Fatalf("failed to read template checkpoint (%s/%s): %v", templateAlgorithmName, templateCheckpointID, err)
	}
	if template == nil {
		t.Fatalf("template checkpoint (%s/%s) not found in database", templateAlgorithmName, templateCheckpointID)
	}

	cp := template.DeepCopy()
	newID := uuid.New().String()
	cp.Id = newID
	cp.Tag = fmt.Sprintf("smoke_%s", newID)
	cp.LifecycleStage = coremodels.LifecycleStageRunning
	cp.ResultUri = ""
	cp.AlgorithmFailureCause = ""
	cp.AlgorithmFailureDetails = ""
	cp.ReceivedAt = time.Now()
	cp.LastModified = time.Now()

	if err := cqlStore.UpsertCheckpoint(cp); err != nil {
		t.Fatalf("failed to upsert provisioned checkpoint %s: %v", newID, err)
	}

	return cp
}

func TestSmoke_CheckRun_InitialState(t *testing.T) {
	cqlStore := newTestStore(t)
	cp := provisionCheckpoint(t, cqlStore)

	url := fmt.Sprintf("%s/algorithm/v1/check/%s/requests/%s", baseURL, cp.Algorithm, cp.Id)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to perform GET %s: %v", url, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 200 OK, got %d. body: %s", resp.StatusCode, string(body))
	}

	var checkResp models.CheckRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		t.Fatalf("failed to decode check run response: %v", err)
	}

	// For a newly provisioned RUNNING checkpoint, is_processed should be false
	if checkResp.IsProcessed {
		t.Errorf("expected is_processed to be false for initial RUNNING checkpoint %s, got true", cp.Id)
	}
}

func TestSmoke_CheckRun_NotFound(t *testing.T) {
	algorithmName := "non-existent-algorithm"
	requestId := uuid.New().String()

	url := fmt.Sprintf("%s/algorithm/v1/check/%s/requests/%s", baseURL, algorithmName, requestId)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to perform GET %s: %v", url, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 404 Not Found, got %d. body: %s", resp.StatusCode, string(body))
	}
}

func TestSmoke_CompleteRun_And_CheckRun(t *testing.T) {
	cqlStore := newTestStore(t)
	cp := provisionCheckpoint(t, cqlStore)

	payload := models.AlgorithmResult{
		ResultUri: fmt.Sprintf("http://localhost:9000/smoke-test-result-%s", cp.Id),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal complete payload: %v", err)
	}

	completeURL := fmt.Sprintf("%s/algorithm/v1/complete/%s/requests/%s", baseURL, cp.Algorithm, cp.Id)
	resp, err := http.Post(completeURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("failed to perform POST %s: %v", completeURL, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 202 Accepted, got %d. body: %s", resp.StatusCode, string(body))
	}

	// Completion actor processes asynchronously in the background.
	// Poll check endpoint until is_processed becomes true or timeout occurs.
	checkURL := fmt.Sprintf("%s/algorithm/v1/check/%s/requests/%s", baseURL, cp.Algorithm, cp.Id)
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
			_ = checkRespHttp.Body.Close()
			continue
		}

		var checkResult models.CheckRunResponse
		decodeErr := json.NewDecoder(checkRespHttp.Body).Decode(&checkResult)
		_ = checkRespHttp.Body.Close()

		if decodeErr == nil && checkResult.IsProcessed {
			completed = true
			break
		}
	}

	if !completed {
		t.Fatalf("timed out waiting for checkpoint %s to be marked as processed", cp.Id)
	}
}

func TestSmoke_CompleteRun_InvalidPayload(t *testing.T) {
	cqlStore := newTestStore(t)
	cp := provisionCheckpoint(t, cqlStore)

	completeURL := fmt.Sprintf("%s/algorithm/v1/complete/%s/requests/%s", baseURL, cp.Algorithm, cp.Id)
	resp, err := http.Post(completeURL, "application/json", strings.NewReader("invalid-json"))
	if err != nil {
		t.Fatalf("failed to perform POST %s: %v", completeURL, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 400 Bad Request for malformed payload, got %d. body: %s", resp.StatusCode, string(body))
	}
}
