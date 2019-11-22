package parser

import (
	"fmt"
	"runtime"

	"github.com/yoheimuta/go-protoparser/issues"
)

func (p *Parser) unexpected(expected string) issues.ParseError {
	_, file, line, _ := runtime.Caller(1)

	pe := issues.ParseError{
		Filename: p.lex.Pos.Filename,
		Line:     p.lex.Pos.Line,
		Column:   p.lex.Pos.Column,
		Found:    p.lex.String(),
		Expected: expected,
	}
	pe.SetOccured(file, line)
	return pe
}

func (p *Parser) unexpectedf(
	format string,
	a ...interface{},
) error {
	return p.unexpected(fmt.Sprintf(format, a...))
}
