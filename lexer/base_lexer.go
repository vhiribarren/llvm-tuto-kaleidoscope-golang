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
	"fmt"
	"unicode/utf8"
)

type lexerErrorType int

const (
	LexerErrorEOF lexerErrorType = iota
	LexerErrorBadRune
)

func (e lexerErrorType) String() string {
	switch e {
	case LexerErrorEOF:
		return "LexerErrorEOF"
	case LexerErrorBadRune:
		return "LexerErrorBadRune"
	default:
		panic("Unknown error type")
	}
}

type LexerError struct {
	errorType lexerErrorType
	message   string
}

func (l LexerError) Error() string {
	return fmt.Sprintf("%s: %s", l.errorType, l.message)
}

func newErrorEOF() error {
	return LexerError{errorType: LexerErrorEOF, message: "End of buffer reached"}
}

func newErrorBadRune(message string, args ...interface{}) error {
	return LexerError{errorType: LexerErrorBadRune, message: fmt.Sprintf(message, args...)}
}

type BaseLexer struct {
	pos    int
	buffer string
}

func NewBaseLexer(data string) BaseLexer {
	return BaseLexer{pos: 0, buffer: data}
}

func (l BaseLexer) PeekNext() (rune, error) {
	if l.pos >= len(l.buffer) {
		return 0, newErrorEOF()
	}
	val, _ := utf8.DecodeRuneInString(l.buffer[l.pos:])
	return val, nil
}

func (l *BaseLexer) ConsumeNext() (rune, error) {
	if l.pos >= len(l.buffer) {
		return 0, newErrorEOF()
	}
	val, width := utf8.DecodeRuneInString(l.buffer[l.pos:])
	l.pos += width
	return val, nil
}

func (l *BaseLexer) ConsumeRune(val rune) error {
	candidate, err := l.ConsumeNext()
	if err != nil {
		return err
	}
	if candidate == val {
		return nil
	}
	return newErrorBadRune("Looked for: %c but found: %c", val, candidate)
}

func (l *BaseLexer) ConsumeString(val string) error {
	for _, targetRune := range val {
		candidateRune, err := l.ConsumeNext()
		if err != nil {
			return err
		}
		if candidateRune != targetRune {
			return newErrorBadRune("While searching for: %s, looked for: %c but found: %c", val, targetRune, candidateRune)
		}
	}
	return nil
}

func (l *BaseLexer) ConsumeWhitespaces() {
	for {
		val, err := l.PeekNext()
		if err != nil {
			return
		}
		if !IsWhitespace(val) {
			return
		}
		l.ConsumeNext()
	}
}

func IsWhitespace(val rune) bool {
	switch val {
	case ' ', '\t', '\n', '\f', '\r':
		return true
	}
	return false
}
