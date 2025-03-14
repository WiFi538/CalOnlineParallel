package calculator

import (
	"errors"
	"strconv"
	"unicode"
)

func Calc(expression string) (float64, error) {
	numStack := make([]float64, 0)
	opStack := make([]rune, 0)

	doOp := func(op rune) error {
		if len(numStack) < 2 {
			return errors.New("недостаточно операндов для операции")
		}
		b := numStack[len(numStack)-1]
		a := numStack[len(numStack)-2]
		numStack = numStack[:len(numStack)-2]
		var result float64
		switch op {
		case '+':
			result = a + b
		case '-':
			result = a - b
		case '*':
			result = a * b
		case '/':
			if b == 0 {
				return errors.New("деление на ноль")
			}
			result = a / b
		default:
			return errors.New("неизвестная операция")
		}
		numStack = append(numStack, result)
		return nil
	}

	priority := func(op rune) int {
		switch op {
		case '+', '-':
			return 1
		case '*', '/':
			return 2
		}
		return 0
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])
		if unicode.IsDigit(char) {
			num, err := strconv.ParseFloat(string(char), 64)
			if err != nil {
				return 0, err
			}
			numStack = append(numStack, num)
		} else if char == '(' {
			opStack = append(opStack, char)
		} else if char == ')' {
			for len(opStack) > 0 && opStack[len(opStack)-1] != '(' {
				if err := doOp(opStack[len(opStack)-1]); err != nil {
					return 0, err
				}
				opStack = opStack[:len(opStack)-1]
			}
			if len(opStack) == 0 {
				return 0, errors.New("несогласованные скобки")
			}
			opStack = opStack[:len(opStack)-1]
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			for len(opStack) > 0 && priority(opStack[len(opStack)-1]) >= priority(char) {
				if err := doOp(opStack[len(opStack)-1]); err != nil {
					return 0, err
				}
				opStack = opStack[:len(opStack)-1]
			}
			opStack = append(opStack, char)
		} else if char != ' ' {
			return 0, errors.New("недопустимый символ в выражении")
		}
	}

	for len(opStack) > 0 {
		if opStack[len(opStack)-1] == '(' {
			return 0, errors.New("непарные скобки")
		}
		if err := doOp(opStack[len(opStack)-1]); err != nil {
			return 0, err
		}
		opStack = opStack[:len(opStack)-1]
	}

	if len(numStack) != 1 {
		return 0, errors.New("некорректное выражение")
	}

	return numStack[0], nil
}
