package lexer

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/parser/meta"
)

func (lex *Lexer) unexpected(found, expected string) error {
	err := &meta.Error{
		Pos:      lex.Pos.Position,
		Expected: expected,
		Found:    lex.Text,
	}
	if lex.debug {
		_, file, line, _ := runtime.Caller(1)
		err.SetOccured(file, line)
	}
	return err
}
