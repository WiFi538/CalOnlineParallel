package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
	var errorLogged bool
	for {
		task, err := getTask()
		if err != nil {
			if !errorLogged {
				log.Printf("Error getting task: %v", err)
				errorLogged = true
			}
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

		// Сброс флага после успешного выполнения задачи
		errorLogged = false
	}
}

func getTask() (*Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var taskResponse struct {
		Task *Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&taskResponse); err != nil {
		return nil, err
	}

	if taskResponse.Task == nil {
		return nil, errors.New("no task available")
	}

	return taskResponse.Task, nil
}

func calculateTask(task *Task) (float64, error) {
	arg1, err := strconv.ParseFloat(task.Arg1, 64)
	if err != nil {
		return 0, err
	}
	arg2, err := strconv.ParseFloat(task.Arg2, 64)
	if err != nil {
		return 0, err
	}

	switch task.Operation {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, errors.New("деление на ноль")
		}
		return arg1 / arg2, nil
	default:
		return 0, errors.New("неизвестная операция")
	}
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
