package orchestrator

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Orchestrator struct {
	expressions map[string]*Expression
	tasks       []*Task
	mu          sync.Mutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		expressions: make(map[string]*Expression),
		tasks:       make([]*Task, 0),
	}
}

func (o *Orchestrator) AddExpression(expression string) (string, error) {
	id := uuid.New().String()
	expr := &Expression{
		ID:     id,
		Status: "pending",
	}

	o.mu.Lock()
	o.expressions[id] = expr
	o.mu.Unlock()

	tasks, err := ParseExpression(expression)
	if err != nil {
		return "", err
	}

	o.mu.Lock()
	o.tasks = append(o.tasks, tasks...)
	o.mu.Unlock()

	return id, nil
}

func (o *Orchestrator) GetExpressions() []*Expression {
	o.mu.Lock()
	defer o.mu.Unlock()

	expressions := make([]*Expression, 0, len(o.expressions))
	for _, expr := range o.expressions {
		expressions = append(expressions, expr)
	}
	return expressions
}

func (o *Orchestrator) GetExpressionByID(id string) (*Expression, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	expr, exists := o.expressions[id]
	if !exists {
		return nil, errors.New("expression not found")
	}
	return expr, nil
}

func (o *Orchestrator) GetTask() (*Task, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.tasks) == 0 {
		return nil, errors.New("no tasks available")
	}

	task := o.tasks[0]
	o.tasks = o.tasks[1:]
	return task, nil
}

func (o *Orchestrator) CompleteTask(taskID string, result float64) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, expr := range o.expressions {
		if expr.Status == "pending" {
			expr.Result = result
			expr.Status = "completed"
			return nil
		}
	}
	return errors.New("task not found")
}
