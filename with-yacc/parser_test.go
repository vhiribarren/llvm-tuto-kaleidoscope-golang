package main

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
		Parse(input)
	}
}

func TestInvalidInput(t *testing.T) {
	invalidInputs := [...]string{
		"def a b c",
		"extern extern",
	}
	currentIndex := len(invalidInputs) - 1
	var catchPanic func()
	catchPanic = func() {
		if r := recover(); r == nil {
			t.Fatal("The code did not panic")
		}
		currentIndex--
		if currentIndex < 0 {
			return
		}
		defer catchPanic()
		Parse(invalidInputs[currentIndex])
	}
	defer catchPanic()
	Parse(invalidInputs[currentIndex])
}
