package lexer

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/issues"
)

func (lex *Lexer) unexpected(found, expected string) error {
	pe := issues.ParseError{
		Filename: lex.Pos.Filename,
		Line:     lex.Pos.Line,
		Column:   lex.Pos.Column,
		Found:    found,
		Expected: expected,
	}
	if lex.debug {
		_, file, line, _ := runtime.Caller(1)
		pe.SetOccured(file, line)
	}

	return pe
}
