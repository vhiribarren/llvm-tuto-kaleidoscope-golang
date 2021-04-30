package lexer

import (
	"testing"
)

func TestPeekAndConsume(t *testing.T) {
	lexer := NewBaseLexer("abc")
	peekedResult, err := lexer.PeekNext()
	if err != nil {
		t.Fatal(err)
	}
	consumedResult, err := lexer.ConsumeNext()
	if err != nil {
		t.Fatal(err)
	}
	if peekedResult != consumedResult {
		t.Fatal("Peeked result and result are not the same")
	}
	newResult, err := lexer.ConsumeNext()
	if newResult == consumedResult {
		t.Fatal("The test was set up so a new consumption would be different than the previous")
	}
}

func TestConsumeValidString(t *testing.T) {
	lexer := NewBaseLexer("hello")
	if err := lexer.ConsumeString("hello"); err != nil {
		t.Fatal(err)
	}
}

func TestConsumeInvalidString(t *testing.T) {
	lexer := NewBaseLexer("hello")
	if err := lexer.ConsumeString("world"); err == nil {
		t.Fatal("Was waiting for an error")
	}
}

func TestWhitespaceConsumptionSimple(t *testing.T) {
	lexer := NewBaseLexer("   \f   \n  \t")
	lexer.ConsumeWhitespaces()
	if _, err := lexer.ConsumeNext(); err == nil {
		t.Fatal("Was waiting for an error")
	}
}

func TestWhitespaceConsumptionWithWords(t *testing.T) {
	lexer := NewBaseLexer("   hello    world")
	lexer.ConsumeWhitespaces()
	if err := lexer.ConsumeString("hello"); err != nil {
		t.Fatal(err)
	}
	lexer.ConsumeWhitespaces()
	if err := lexer.ConsumeString("world"); err != nil {
		t.Fatal(err)
	}
}
