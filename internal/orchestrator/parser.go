package orchestrator

import (
	"errors"
	"strconv"
	"strings"
)

func ParseExpression(expression string) ([]*Task, error) {
	tokens := tokenize(expression)
	tasks, err := shuntingYard(tokens)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func tokenize(expression string) []string {
	tokens := make([]string, 0)
	var currentToken strings.Builder

	for _, char := range expression {
		if char == ' ' {
			continue
		}
		if char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			currentToken.WriteRune(char)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func shuntingYard(tokens []string) ([]*Task, error) {
	opStack := make([]string, 0)
	output := make([]*Task, 0)
	numStack := make([]float64, 0)

	for _, token := range tokens {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
			numStack = append(numStack, num)
		} else if token == "(" {
			opStack = append(opStack, token)
		} else if token == ")" {
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				op := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if len(numStack) < 2 {
					return nil, errors.New("недостаточно операндов для операции")
				}
				arg2 := numStack[len(numStack)-1]
				arg1 := numStack[len(numStack)-2]
				numStack = numStack[:len(numStack)-2]
				output = append(output, &Task{Arg1: strconv.FormatFloat(arg1, 'f', -1, 64), Arg2: strconv.FormatFloat(arg2, 'f', -1, 64), Operation: op})
				result := calculate(arg1, arg2, op)
				numStack = append(numStack, result)
			}
			if len(opStack) == 0 {
				return nil, errors.New("unbalanced parentheses")
			}
			opStack = opStack[:len(opStack)-1]
		} else if isOperator(token) {
			for len(opStack) > 0 && precedence(opStack[len(opStack)-1]) >= precedence(token) {
				op := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if len(numStack) < 2 {
					return nil, errors.New("недостаточно операндов для операции")
				}
				arg2 := numStack[len(numStack)-1]
				arg1 := numStack[len(numStack)-2]
				numStack = numStack[:len(numStack)-2]
				output = append(output, &Task{Arg1: strconv.FormatFloat(arg1, 'f', -1, 64), Arg2: strconv.FormatFloat(arg2, 'f', -1, 64), Operation: op})
				result := calculate(arg1, arg2, op)
				numStack = append(numStack, result)
			}
			opStack = append(opStack, token)
		} else {
			return nil, errors.New("invalid token")
		}
	}

	for len(opStack) > 0 {
		if opStack[len(opStack)-1] == "(" {
			return nil, errors.New("unbalanced parentheses")
		}
		op := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]
		if len(numStack) < 2 {
			return nil, errors.New("недостаточно операндов для операции")
		}
		arg2 := numStack[len(numStack)-1]
		arg1 := numStack[len(numStack)-2]
		numStack = numStack[:len(numStack)-2]
		output = append(output, &Task{Arg1: strconv.FormatFloat(arg1, 'f', -1, 64), Arg2: strconv.FormatFloat(arg2, 'f', -1, 64), Operation: op})
		result := calculate(arg1, arg2, op)
		numStack = append(numStack, result)
	}

	return output, nil
}

func calculate(arg1, arg2 float64, op string) float64 {
	switch op {
	case "+":
		return arg1 + arg2
	case "-":
		return arg1 - arg2
	case "*":
		return arg1 * arg2
	case "/":
		return arg1 / arg2
	default:
		return 0
	}
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}
