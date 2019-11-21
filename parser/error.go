package parser

import (
	"fmt"
	"runtime"

	"github.com/yoheimuta/go-protoparser/issues"
)

func (p *Parser) unexpected(expected string) issues.ParseError {
	_, file, line, _ := runtime.Caller(1)

	return issues.NewParseError(
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
