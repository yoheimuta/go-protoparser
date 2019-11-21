package lexer

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/issues"
)

func (lex *Lexer) unexpected(found, expected string) error {
	file := ""
	line := 0
	if lex.debug {
		_, file, line, _ = runtime.Caller(1)
	}
	return issues.NewParseError(found, expected, file, line)
}
