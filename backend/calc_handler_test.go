package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalcHandler_Options(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/api/calc", nil)
	rr := httptest.NewRecorder()
	calcHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestCalcHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/calc", nil)
	rr := httptest.NewRecorder()
	calcHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestCalcHandler_BadRequests(t *testing.T) {
	tests := []struct {
		name string
		body any
		want int
	}{
		{"invalid json", "not-json", http.StatusBadRequest},
		{"missing op/a", map[string]any{}, http.StatusBadRequest},
		{"missing b for add", map[string]any{"op": "add", "a": 1}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}
			req := httptest.NewRequest(http.MethodPost, "/api/calc", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			calcHandler(rr, req)

			if rr.Code != tt.want {
				t.Fatalf("expected status %d, got %d; body=%s", tt.want, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestCalcHandler_SuccessAndErrors(t *testing.T) {
	t.Run("add success", func(t *testing.T) {
		reqBody := map[string]any{"op": "add", "a": 11, "b": 12}
		b, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/calc", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		calcHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d; body=%s", rr.Code, rr.Body.String())
		}
		var resp CalcResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if resp.Error != nil {
			t.Fatalf("expected no error, got %v", *resp.Error)
		}
		if resp.Result == nil || *resp.Result != 23 {
			t.Fatalf("expected result 23, got %+v", resp.Result)
		}
	})

	t.Run("div by zero", func(t *testing.T) {
		reqBody := map[string]any{"op": "div", "a": 1, "b": 0}
		b, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/calc", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		calcHandler(rr, req)

		// handler returns 400 with JSON error on Compute errors
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d; body=%s", rr.Code, rr.Body.String())
		}
		var resp CalcResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if resp.Error == nil {
			t.Fatalf("expected error message, got nil")
		}
	})
}
