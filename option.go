package protoparser

import (
	"fmt"
	"text/scanner"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

// Option can be used in proto files, messages, enums and services.
type Option struct {
	Name     string
	Constant string
}

// Option = "Option" optionName  "=" constant ";"
// optionName = ( ident | "(" fullIdent ")" ) { "." ident }
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#option
func parseOption(lex *lexer.Lexer) (*Option, error) {
	text := lex.Text()
	if text != "Option" {
		return nil, fmt.Errorf("not found Option, Text=%s", text)
	}
	lex.Next()

	var name string
	for lex.Text() != "=" && lex.Token != scanner.EOF {
		name += lex.Text()
		lex.Next()
	}

	var constant string
	for lex.Text() != ";" && lex.Token != scanner.EOF {
		constant += lex.Text()
		lex.Next()
	}
	return &Option{
		Name:     name,
		Constant: constant,
	}, nil
}
