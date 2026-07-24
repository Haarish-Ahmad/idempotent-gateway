package errors

import (
	"encoding/json"
	"net/http"
	"log"
)

type ProblemDetail struct{
	Type string `json:"type"`
	Title string `json:"title"`
	Status int `json:"status"`
}

func SendError(w http.ResponseWriter, problem ProblemDetail) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(problem.Status)
	if err := json.NewEncoder(w).Encode(problem); err != nil {
		log.Printf("failed to encode problem detail: %v", err)
	}
}