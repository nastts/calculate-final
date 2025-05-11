package calculate_test

import (
	"testing"

	"github.com/nastts/final-calculator/calculate"
)

func TestCalc(t *testing.T){
	tests := []struct{
		expression string
		result float64
	}{
		{"1+1", 2},
		{"3+3*6", 21},
		{"1+8/2*4", 17},
		{"(1+1)*2", 4},
		{"10-2+3", 11},
		{"5*(2+3)", 25},
		{"(8/4)+(3*2)", 8},
		{"7-(3+2)", 2},
		{"6/2*(1+2)", 9},
		{"(3+5)*(2-1)", 8},
		{"(10-3)*(5+ 2)", 49},
		{"(6+2)*(3/2)", 12},
		{"(5 + 5) / (10 / 2)", 2},
		{"(8 + 4) * (2 - 1)", 12},
		
	}
	for _, testCase := range tests {
		
		result, _ := calculate.Calc(testCase.expression)
		if result != testCase.result {
			t.Errorf("%s=%v, want %v", testCase.expression, result, testCase.result)
		}
	}
}

