package parser

import (
	"fmt"
	"runtime"

	"github.com/yoheimuta/go-protoparser/errors"
)

func (p *Parser) unexpected(expected string) errors.ParseError {
	_, file, line, _ := runtime.Caller(1)

	return errors.NewParseError(
		p.lex.String(),
		expected,
		file,
		line,
	)
}

func (p *Parser) unexpectedf(
	format string,
	a ...interface{},
) error {
	return p.unexpected(fmt.Sprintf(format, a...))
}
