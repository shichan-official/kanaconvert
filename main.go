package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type SumRequest struct {
	Num1 float64 `json:"num1"`
	Num2 float64 `json:"num2"`
}

type SumResponse struct {
	Sum float64 `json:"sum"`
}

// CORS Middleware
func enableCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, or specify your frontend URL
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	}
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
	http.HandleFunc("/sum", enableCors(sumHandler))

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
