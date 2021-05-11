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

func TestValidMedley(t *testing.T) {
	input := "machin123    def   defextern <= 123  extern 456hello  #comment def"
	targetResults := []KaleidoTokenContext{
		{KTokenIdentifier, "machin123"},
		{KTokenDef, ""},
		{KTokenIdentifier, "defextern"},
		{KTokenSymbol, "<"},
		{KTokenSymbol, "="},
		{KTokenNumber, "123"},
		{KTokenExtern, ""},
		{KTokenNumber, "456"},
		{KTokenIdentifier, "hello"},
		{KTokenEOF, ""},
	}
	lexer := NewKaleidoLexer(input)
	for i := 0; i < len(targetResults); i++ {
		result := lexer.NextToken()
		if result.Token != targetResults[i].Token || result.Value != targetResults[i].Value {
			t.Fatalf("Was waiting for: %v but creceived: %v", result, &targetResults[i])
		}
	}

}
