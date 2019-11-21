package parser

import (
	"fmt"
	"runtime"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

type ParseError struct {
	Lexer *lexer.Lexer

	Expected  string
	occuredIn string
	occuredAt int
}

func (pe ParseError) String() string {
	return fmt.Sprintf("found %q(Token=%v, Pos=%s) but expected [%s] at %s:%d", pe.Lexer.Text, pe.Lexer.Token, pe.Lexer.Pos, pe.Expected, pe.occuredIn, pe.occuredAt)
}

func (pe ParseError) Error() string {
	return pe.String()
}

func (p *Parser) unexpected(expected string) ParseError {
	_, file, line, _ := runtime.Caller(1)

	return ParseError{
		Lexer:     p.lex,
		Expected:  expected,
		occuredIn: file,
		occuredAt: line,
	}
}

func (p *Parser) unexpectedf(
	format string,
	a ...interface{},
) error {
	return p.unexpected(fmt.Sprintf(format, a...))
}
