package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Request body structure
type SumRequest struct {
	Num1 float64 `json:"num1"`
	Num2 float64 `json:"num2"`
}

// Response body structure
type SumResponse struct {
	Sum float64 `json:"sum"`
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	var req SumRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sum := req.Num1 + req.Num2
	res := SumResponse{Sum: sum}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/sum", sumHandler)

	port := "8080" // Default port
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
