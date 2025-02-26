package main

import (
	"log"
	"net/http"

	"github.com/wifi538/CalOnlineParallel/internal/orchestrator"
)

func main() {
	http.HandleFunc("/api/v1/calculate", orchestrator.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", orchestrator.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/:id", orchestrator.HandleGetExpressionByID)
	http.HandleFunc("/internal/task", orchestrator.HandleGetTask)
	http.HandleFunc("/internal/task", orchestrator.HandlePostTaskResult)

	log.Println("Orchestrator is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
