package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wifi538/CalOnlineParallel/internal/orchestrator"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", orchestrator.HandleCalculate).Methods("POST")
	r.HandleFunc("/api/v1/expressions", orchestrator.HandleGetExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", orchestrator.HandleGetExpressionByID).Methods("GET")
	r.HandleFunc("/internal/task", orchestrator.HandleGetTask).Methods("GET")
	r.HandleFunc("/internal/task", orchestrator.HandlePostTaskResult).Methods("POST")

	log.Println("Orchestrator is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
