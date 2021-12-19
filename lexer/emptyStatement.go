package lexer

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
)

// ReadEmptyStatement reads an emptyStatement.
//  emptyStatement = ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#emptystatement
func (lex *Lexer) ReadEmptyStatement() error {
	lex.Next()

	if lex.Token == scanner.TSEMICOLON {
		return nil
	}
	lex.UnNext()
	return lex.unexpected(lex.Text, ";")
}
