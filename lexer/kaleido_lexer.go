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
	"strings"
)

type KaleidoToken int

//go:generate go run golang.org/x/tools/cmd/stringer -type=KaleidoToken
const (
	KTokenEOF KaleidoToken = iota
	KTokenDef
	KTokenExtern
	KTokenIdentifier
	KTokenNumber
	KTokenSymbol
)

type KaleidoTokenContext struct {
	Token KaleidoToken
	Value string
}

func emitEOF() *KaleidoTokenContext {
	return &KaleidoTokenContext{Token: KTokenEOF, Value: ""}
}

func emitDef() *KaleidoTokenContext {
	return &KaleidoTokenContext{Token: KTokenDef, Value: ""}
}

func emitExtern() *KaleidoTokenContext {
	return &KaleidoTokenContext{Token: KTokenExtern, Value: ""}
}

func emitIdentifier(identifier string) *KaleidoTokenContext {
	return &KaleidoTokenContext{Token: KTokenIdentifier, Value: identifier}
}

func emitNumber(val string) *KaleidoTokenContext {
	return &KaleidoTokenContext{Token: KTokenNumber, Value: val}
}

func emitSymbol(val rune) *KaleidoTokenContext {
	return &KaleidoTokenContext{Token: KTokenSymbol, Value: string(val)}
}

type KaleidoLexer struct {
	BaseLexer
}

func NewKaleidoLexer(data string) KaleidoLexer {
	return KaleidoLexer{BaseLexer: NewBaseLexer(data)}
}

func (l *KaleidoLexer) NextToken() *KaleidoTokenContext {
	for {
		l.ConsumeWhitespaces()
		val, err := l.PeekNext()
		switch {
		case err != nil:
			return emitEOF()
		case isAlphabetic(val):
			result := l.consumeGreedAlphanum()
			switch result {
			case "def":
				return emitDef()
			case "extern":
				return emitExtern()
			default:
				return emitIdentifier(result)
			}
		case isNumeral(val):
			result := l.consumeGreedNumber()
			return emitNumber(result)
		case val == '#':
			l.consumeGreedCommentLine()
		default:
			l.ConsumeNext()
			return emitSymbol(val)
		}

	}
}

func (l *KaleidoLexer) consumeGreedCommentLine() {
	for {
		val, err := l.PeekNext()
		if err != nil {
			return
		}
		if val == '\n' {
			return
		}
		l.ConsumeNext()
	}
}

func (l *KaleidoLexer) consumeGreedAlphanum() string {
	var builder strings.Builder
	for {
		val, err := l.PeekNext()
		if err != nil {
			return builder.String()
		}
		if !isAlphanum(val) {
			return builder.String()
		}
		char, _ := l.ConsumeNext()
		builder.WriteRune(char)
	}
}

func (l *KaleidoLexer) consumeGreedNumber() string {
	var builder strings.Builder
	for {
		val, err := l.PeekNext()
		if err != nil {
			return builder.String()
		}
		if !isDigit(val) {
			return builder.String()
		}
		char, _ := l.ConsumeNext()
		builder.WriteRune(char)
	}
}

func isAlphanum(val rune) bool {
	if isNumeral(val) || isAlphabetic(val) {
		return true
	}
	return false
}

func isAlphabetic(val rune) bool {
	if (val >= 'a' && val <= 'z') || (val >= 'A' && val <= 'Z') {
		return true
	}
	return false
}

func isDigit(val rune) bool {
	// As written in the LLVM Kaleido tutorial, it can accept invalid numbers for now...
	if isNumeral(val) || val == '.' {
		return true
	}
	return false
}

func isNumeral(val rune) bool {
	if val >= '0' && val <= '9' {
		return true
	}
	return false
}
