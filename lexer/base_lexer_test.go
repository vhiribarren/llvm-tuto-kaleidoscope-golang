/*
MIT License

Copyright (c) 2021 Vincent Hiribarren

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
