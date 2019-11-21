package lexer

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/errors"
)

func (lex *Lexer) unexpected(found, expected string) error {
	file := ""
	line := 0
	if lex.debug {
		_, file, line, _ = runtime.Caller(1)
	}
	return errors.NewParseError(found, expected, file, line)
}
