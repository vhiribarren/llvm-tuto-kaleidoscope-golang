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

func (l *BaseLexer) ConsumeWhitespaces() error {
	for {
		val, err := l.PeekNext()
		if err != nil {
			return nil
		}
		if !isWhitespace(val) {
			return nil
		}
		l.ConsumeNext()
	}
}

func isWhitespace(val rune) bool {
	switch val {
	case ' ', '\t', '\n', '\f', '\r':
		return true
	}
	return false
}
