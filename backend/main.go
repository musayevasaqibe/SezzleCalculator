package main

import (
	"encoding/json"
	"log"
	"net/http"
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

func calcHandler(w http.ResponseWriter, r *http.Request) {
	// Allow CORS for local dev
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// simple validation
	if req.Op == "" || req.A == nil {
		http.Error(w, "missing op or a", http.StatusBadRequest)
		return
	}

	// For operations that require b, ensure it's present
	needsB := req.Op == "add" || req.Op == "sub" || req.Op == "mul" || req.Op == "div" || req.Op == "pow" || req.Op == "pct"
	if needsB && req.B == nil {
		http.Error(w, "missing b for this operation", http.StatusBadRequest)
		return
	}

	result, err := Compute(req.Op, *req.A, func() float64 {
		if req.B == nil {
			return 0
		}
		return *req.B
	}())

	var resp CalcResponse
	if err != nil {
		msg := err.Error()
		resp.Error = &msg
		resp.Result = nil
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	resp.Result = &result
	resp.Error = nil
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
