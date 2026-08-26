package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestEvalEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := evalRequest{Expr: "x + 2 * y", Vars: map[string]float64{"x": 3, "y": 4}}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/eval", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp evalResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Result != 11 {
		t.Errorf("expected 11, got %f", resp.Result)
	}
}

func TestEvalEndpoint_Empty(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	body := []byte(`{"expr":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/eval", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestEvalEndpoint_ParseError(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := evalRequest{Expr: "1 + * 2"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/eval", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestValidateEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := validateRequest{Expr: "x + y", Vars: []string{"x", "y"}}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/validate", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp validateResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.OK {
		t.Errorf("expected valid, got errors: %v", resp.Errors)
	}
}

func TestValidateEndpoint_UnknownVar(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := validateRequest{Expr: "x + z", Vars: []string{"x"}}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/validate", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp validateResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.OK {
		t.Error("expected validation failure for unknown var z")
	}
}

func TestCompileEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := compileRequest{Expr: "2 + 3 * 4", Vars: map[string]float64{}}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/compile", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp compileResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Result != 14 {
		t.Errorf("expected 14, got %f", resp.Result)
	}
	if resp.Instructions <= 0 {
		t.Error("expected positive instruction count")
	}
}

func TestBatchEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := batchRequest{
		Expressions: []string{"x + 1", "x * 2", "1 * * BAD"},
		Vars:        map[string]float64{"x": 5},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/batch", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Results []batchResult `json:"results"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(resp.Results))
	}
	if resp.Results[0].Result == nil || *resp.Results[0].Result != 6 {
		t.Errorf("expected first result=6")
	}
	if resp.Results[2].Error == "" {
		t.Error("expected error for invalid expression")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	endpoints := []string{"/api/eval", "/api/validate", "/api/compile", "/api/batch"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", ep, rec.Code)
		}
	}
}
