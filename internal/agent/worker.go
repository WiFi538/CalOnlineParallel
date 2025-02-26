package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wifi538/CalOnlineParallel/pkg/calculator"
)

type Task struct {
	ID            string `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

func Worker() {
	for {
		task, err := getTask()
		if err != nil {
			log.Printf("Error getting task: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		result, err := calculateTask(task)
		if err != nil {
			log.Printf("Error calculating task: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = sendResult(task.ID, result)
		if err != nil {
			log.Printf("Error sending result: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

func getTask() (*Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func calculateTask(task *Task) (float64, error) {
	expression := task.Arg1 + " " + task.Operation + " " + task.Arg2
	return calculator.Calc(expression)
}

func sendResult(taskID string, result float64) error {
	taskResult := TaskResult{
		ID:     taskID,
		Result: result,
	}

	data, err := json.Marshal(taskResult)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
