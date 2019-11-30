package parser

import (
	"fmt"
	"runtime"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func (p *Parser) unexpected(expected string) error {
	_, file, line, _ := runtime.Caller(1)
	err := &meta.Error{
		Pos:      p.lex.Pos.Position,
		Expected: expected,
		Found:    fmt.Sprintf("%q(Token=%v, Pos=%s)", p.lex.Text, p.lex.Token, p.lex.Pos),
	}
	err.SetOccured(file, line)
	return err
}

func (p *Parser) unexpectedf(
	format string,
	a ...interface{},
) error {
	return p.unexpected(fmt.Sprintf(format, a...))
}
