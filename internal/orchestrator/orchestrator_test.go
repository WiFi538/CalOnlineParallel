package orchestrator

import (
	"testing"
)

func TestAddExpression(t *testing.T) {
	o := NewOrchestrator()
	id, err := o.AddExpression("2+2*2")
	if err != nil {
		t.Fatalf("Failed to add expression: %v", err)
	}

	if id == "" {
		t.Fatalf("Expected non-empty ID")
	}
}

func TestGetExpressions(t *testing.T) {
	o := NewOrchestrator()
	_, _ = o.AddExpression("2+2*2")
	expressions := o.GetExpressions()

	if len(expressions) != 1 {
		t.Fatalf("Expected 1 expression, got %d", len(expressions))
	}
}

func TestGetExpressionByID(t *testing.T) {
	o := NewOrchestrator()
	id, _ := o.AddExpression("2+2*2")
	expr, err := o.GetExpressionByID(id)

	if err != nil {
		t.Fatalf("Failed to get expression by ID: %v", err)
	}

	if expr.ID != id {
		t.Fatalf("Expected ID %s, got %s", id, expr.ID)
	}
}

func TestGetTask(t *testing.T) {
	o := NewOrchestrator()
	_, _ = o.AddExpression("2+2*2")
	task, err := o.GetTask()

	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task == nil {
		t.Fatalf("Expected non-nil task")
	}
}

func TestCompleteTask(t *testing.T) {
	o := NewOrchestrator()
	id, _ := o.AddExpression("2+2*2")
	task, _ := o.GetTask()
	err := o.CompleteTask(task.ID, 6)

	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	expr, _ := o.GetExpressionByID(id)
	if expr.Status != "completed" {
		t.Fatalf("Expected status 'completed', got %s", expr.Status)
	}

	if expr.Result != 6 {
		t.Fatalf("Expected result 6, got %f", expr.Result)
	}
}
