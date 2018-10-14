package lexer

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// ReadEmptyStatement reads an emptyStatement.
// emptyStatement = ";"
func (lex *Lexer2) ReadEmptyStatement() error {
	lex.scanner.Mode = scanner.ScanIdent
	lex.Next()

	if lex.Token == scanner.TSEMICOLON {
		return nil
	}
	lex.UnNext()
	return lex.unexpected(lex.Text, ";")
}
