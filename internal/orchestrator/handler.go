package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var orchestrator = NewOrchestrator()

func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := orchestrator.AddExpression(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	expressions := orchestrator.GetExpressions()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]*Expression{"expressions": expressions})
}

func HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	expr, err := orchestrator.GetExpressionByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]*Expression{"expression": expr})
}

func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := orchestrator.GetTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]*Task{"task": task})
}

func HandlePostTaskResult(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := orchestrator.CompleteTask(req.ID, req.Result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
