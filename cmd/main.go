package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Haarish-Ahmad/idempotent-gateway/internal/errors"
)

func sendProblem(w http.ResponseWriter, errtype string, title string, status int) {
	errors.SendError(w, errors.ProblemDetail{
		Type: errtype,
		Title: title,
		Status: status,
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Welcome to the API Gateway!")
	})

	mux.HandleFunc("/test-error", func(w http.ResponseWriter, r *http.Request){
		sendProblem(w, "http://gateway.local/errors/bad-request", "Bad Request", http.StatusBadRequest)
	})

	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	sendProblem(w, "http://gateway.local/errors/not-found", "Not Found", http.StatusNotFound)
	//})

	log.Println("Starting API Gateway on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}