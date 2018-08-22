package protoparser

import (
	"fmt"
	"text/scanner"
)

// Option can be used in proto files, messages, enums and services.
type Option struct {
	Name     string
	Constant string
}

// option = "option" optionName  "=" constant ";"
// optionName = ( ident | "(" fullIdent ")" ) { "." ident }
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#option
func parseOption(lex *lexer) (*Option, error) {
	text := lex.text()
	if text != "option" {
		return nil, fmt.Errorf("not found option, text=%s", text)
	}
	lex.next()

	var name string
	for lex.text() != "=" && lex.token != scanner.EOF {
		name += lex.text()
		lex.next()
	}

	var constant string
	for lex.text() != ";" && lex.token != scanner.EOF {
		constant += lex.text()
		lex.next()
	}
	return &Option{
		Name:     name,
		Constant: constant,
	}, nil
}
