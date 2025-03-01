package orchestrator

type Task struct {
	ID            string `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
	ExpressionID  string `json:"expression_id"` // Добавляем связь с выражением
}
