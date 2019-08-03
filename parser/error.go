package parser

import (
	"fmt"
	"runtime"
)

func (p *Parser) unexpected(expected string) error {
	_, file, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(" at %s:%d", file, line)
	return fmt.Errorf("found %q(Token=%v, Pos=%s) but expected [%s]%s", p.lex.Text, p.lex.Token, p.lex.Pos, expected, msg)
}

func (p *Parser) unexpectedf(
	format string,
	a ...interface{},
) error {
	return p.unexpected(fmt.Sprintf(format, a...))
}
