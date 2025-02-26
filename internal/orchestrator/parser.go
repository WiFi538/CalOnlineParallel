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

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, &Task{Arg1: token})
		} else if token == "(" {
			opStack = append(opStack, token)
		} else if token == ")" {
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				output = append(output, &Task{Operation: opStack[len(opStack)-1]})
				opStack = opStack[:len(opStack)-1]
			}
			if len(opStack) == 0 {
				return nil, errors.New("unbalanced parentheses")
			}
			opStack = opStack[:len(opStack)-1]
		} else if isOperator(token) {
			for len(opStack) > 0 && precedence(opStack[len(opStack)-1]) >= precedence(token) {
				output = append(output, &Task{Operation: opStack[len(opStack)-1]})
				opStack = opStack[:len(opStack)-1]
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
		output = append(output, &Task{Operation: opStack[len(opStack)-1]})
		opStack = opStack[:len(opStack)-1]
	}

	return output, nil
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
