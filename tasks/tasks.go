package tasks

import (
	"strconv"
	"strings"

	"github.com/nastts/final-calculator/calculate"
)

func Tokenize(expression string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, char := range expression {
		if char == ' ' {
			continue
		}
		if Operator(string(char)) || char == '(' || char == ')' {
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

func ParseFloat(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func Operator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func precedence(op string) int {
	switch op {
	case "*", "/":
		return 2
	case "+", "-":
		return 1

	default:
		return 0
	}
}

func EvaluateRPN(tokens []string) (float64, error) {
	stack := []float64{}

	for _, token := range tokens {
		if ParseFloat(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if Operator(token) {
			if len(stack) < 2 {
				return 0, calculate.ErrInternalServerError
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, calculate.ErrExpressionIsNotValid
				}
				stack = append(stack, a/b)
			default:
				return 0, calculate.ErrExpressionIsNotValid
			}
		} else {
			return 0, calculate.ErrInternalServerError
		}
	}

	if len(stack) != 1 {
		return 0, calculate.ErrInternalServerError
	}
	return stack[0], nil
}



func Calc(expression string) ([]string, error) {
	tokens := Tokenize(expression)
	if len(tokens) == 0 {
		return nil, calculate.ErrInternalServerError
	}

	output := []string{}
	operatorStack := []string{}

	for _, token := range tokens {
		if ParseFloat(token) {
			output = append(output, token)
		} else if Operator(token) {
			for len(operatorStack) > 0 && precedence(operatorStack[len(operatorStack)-1]) >= precedence(token) {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		} else if token == "(" {
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, calculate.ErrExpressionIsNotValid
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else {
			return nil, calculate.ErrExpressionIsNotValid
		}
	}

	for len(operatorStack) > 0 {
		if operatorStack[len(operatorStack)-1] == "(" {
			return nil,calculate.ErrExpressionIsNotValid
		}
		output = append(output, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return output, nil
}
