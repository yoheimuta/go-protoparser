package protoparser

import (
	"fmt"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

// Oneof consists of oneof fields and a oneof name.
type Oneof struct {
	Comments []string
	Name     string
	Fields   []*Field
}

// "oneof" var '{' OneofContent '}'
func parseOneof(lex *lexer.Lexer) (*Oneof, error) {
	text := lex.Text()
	if text != "oneof" {
		return nil, fmt.Errorf("not found oneof, Text=%s", text)
	}

	// get the name {
	lex.Next()
	name := lex.Text()
	lex.Next()
	// }

	// get the content {
	/// consume '{' {
	lex.Next()
	/// }
	fields, _, _, _, err := parseMessageContent(lex)
	if err != nil {
		return nil, err
	}

	/// consume '}' {
	lex.Next()
	/// }
	// }

	return &Oneof{
		Name:   name,
		Fields: fields,
	}, nil
}
