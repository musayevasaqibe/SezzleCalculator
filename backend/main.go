package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type CalcRequest struct {
	Op string   `json:"op"`
	A  *float64 `json:"a"`
	B  *float64 `json:"b,omitempty"`
}

type CalcResponse struct {
	Result *float64 `json:"result"`
	Error  *string  `json:"error"`
}

func main() {
	http.HandleFunc("/api/calc", calcHandler)

	addr := ":8080"
	log.Printf("Backend running on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	// CORS: allow from configured origin for local dev
	origin := os.Getenv("VITE_API_URL")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, CalcResponse{Result: nil, Error: ptr("method not allowed")})
		return
	}

	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, CalcResponse{Result: nil, Error: ptr("invalid json")})
		return
	}

	if req.Op == "" || req.A == nil {
		writeJSON(w, http.StatusBadRequest, CalcResponse{Result: nil, Error: ptr("missing op or a")})
		return
	}

	needsB := req.Op == "add" || req.Op == "sub" || req.Op == "mul" || req.Op == "div" || req.Op == "pow" || req.Op == "pct"
	if needsB && req.B == nil {
		writeJSON(w, http.StatusBadRequest, CalcResponse{Result: nil, Error: ptr("missing b for this operation")})
		return
	}

	b := 0.0
	if req.B != nil {
		b = *req.B
	}

	result, err := Compute(req.Op, *req.A, b)

	if err != nil {
		msg := err.Error()
		writeJSON(w, http.StatusBadRequest, CalcResponse{Result: nil, Error: &msg})
		return
	}

	writeJSON(w, http.StatusOK, CalcResponse{Result: &result, Error: nil})
}

func ptr(s string) *string { return &s }
