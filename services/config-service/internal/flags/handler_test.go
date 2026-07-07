package flags

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() *http.ServeMux {
	repo := NewMemoryRepository()
	handler := NewHandler(repo)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return mux
}

func TestCreateFlagInvalidJSONReturnsBadRequest(t *testing.T) {
	mux := newTestServer()

	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(`{invalid-json}`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateFlagDuplicateReturnsConflict(t *testing.T) {
	mux := newTestServer()

	body := `{
		"key": "new-checkout",
		"name": "New Checkout",
		"description": "Gradual checkout rollout",
		"enabled": false,
		"rolloutPercentage": 0,
		"targetingRules": []
	}`

	firstReq := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body))
	firstReq.Header.Set("Content-Type", "application/json")
	firstRec := httptest.NewRecorder()

	mux.ServeHTTP(firstRec, firstReq)

	if firstRec.Code != http.StatusCreated {
		t.Fatalf("expected first request status %d, got %d: %s", http.StatusCreated, firstRec.Code, firstRec.Body.String())
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body))
	secondReq.Header.Set("Content-Type", "application/json")
	secondRec := httptest.NewRecorder()

	mux.ServeHTTP(secondRec, secondReq)

	if secondRec.Code != http.StatusConflict {
		t.Fatalf("expected second request status %d, got %d: %s", http.StatusConflict, secondRec.Code, secondRec.Body.String())
	}
}

func TestGetMissingFlagReturnsNotFound(t *testing.T) {
	mux := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/flags/missing-flag", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}

func TestListFlagsReturnsOK(t *testing.T) {
	mux := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestPatchFlagUpdatesFields(t *testing.T) {
	mux := newTestServer()

	createBody := `{
		"key": "new-checkout",
		"name": "New Checkout",
		"description": "Gradual checkout rollout",
		"enabled": false,
		"rolloutPercentage": 0,
		"targetingRules": []
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	mux.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	patchBody := `{
		"enabled": true,
		"rolloutPercentage": 10
	}`

	patchReq := httptest.NewRequest(http.MethodPatch, "/flags/new-checkout", strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchRec := httptest.NewRecorder()

	mux.ServeHTTP(patchRec, patchReq)

	if patchRec.Code != http.StatusOK {
		t.Fatalf("expected patch status %d, got %d: %s", http.StatusOK, patchRec.Code, patchRec.Body.String())
	}
}

func TestPatchMissingFlagReturnsNotFound(t *testing.T) {
	mux := newTestServer()

	body := `{
		"enabled": true
	}`

	req := httptest.NewRequest(http.MethodPatch, "/flags/missing-flag", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}

func TestPatchInvalidRolloutReturnsBadRequest(t *testing.T) {
	mux := newTestServer()

	createBody := `{
		"key": "new-checkout",
		"name": "New Checkout",
		"description": "Gradual checkout rollout",
		"enabled": false,
		"rolloutPercentage": 0,
		"targetingRules": []
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	mux.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	patchBody := `{
		"rolloutPercentage": 101
	}`

	patchReq := httptest.NewRequest(http.MethodPatch, "/flags/new-checkout", strings.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchRec := httptest.NewRecorder()

	mux.ServeHTTP(patchRec, patchReq)

	if patchRec.Code != http.StatusBadRequest {
		t.Fatalf("expected patch status %d, got %d: %s", http.StatusBadRequest, patchRec.Code, patchRec.Body.String())
	}
}

func TestEvaluateMissingFlagReturnsOK(t *testing.T) {
	mux := newTestServer()

	body := `{
		"user": {
			"id": "user_123"
		},
		"defaultValue": false
	}`

	req := httptest.NewRequest(http.MethodPost, "/flags/missing-flag/evaluate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	if !strings.Contains(rec.Body.String(), "FLAG_NOT_FOUND") {
		t.Fatalf("expected response body to include FLAG_NOT_FOUND, got %s", rec.Body.String())
	}
}

func TestEvaluateDisabledFlagReturnsOK(t *testing.T) {
	mux := newTestServer()

	createBody := `{
		"key": "new-checkout",
		"name": "New Checkout",
		"description": "Gradual checkout rollout",
		"enabled": false,
		"rolloutPercentage": 0,
		"targetingRules": []
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	mux.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	body := `{
		"user": {
			"id": "user_123"
		},
		"defaultValue": true
	}`

	req := httptest.NewRequest(http.MethodPost, "/flags/new-checkout/evaluate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	if !strings.Contains(rec.Body.String(), "FLAG_DISABLED") {
		t.Fatalf("expected response body to include FLAG_DISABLED, got %s", rec.Body.String())
	}
}
