package yacc

import (
	"testing"
)

func TestSomeValidInput(t *testing.T) {
	validInputs := [...]string{
		"extern test(a, b, c);",
		"def test(arg) hello",
		"1+1",
	}
	for _, input := range validInputs {
		BuildKaleidoAST(input)
	}
}

func TestInvalidInput(t *testing.T) {
	invalidInputs := [...]string{
		"def a b c",
		"extern extern",
	}
	for _, input := range invalidInputs {
		_, err := BuildKaleidoAST(input)
		if err == nil {
			t.Error("The code did not panic")
		}
	}
}
