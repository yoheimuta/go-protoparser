package protoparser

import (
	"text/scanner"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

// Field is the basic element of a protocol buffer message.
type Field struct {
	Comments    []string
	Type        *Type
	Name        string
	HasRepeated bool
}

// type name = number validator';'
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#normal_field
func parseField(lex *lexer.Lexer) *Field {
	field := &Field{}

	// check HasRepeated {
	text := lex.Text()
	if text == "repeated" {
		field.HasRepeated = true
		lex.Next()
	}
	// }

	for lex.Text() != ";" && lex.Token != scanner.EOF {
		// get the type {
		if field.Type == nil {
			field.Type = parseType(lex)
			continue
		}
		// }

		// get the name {
		token := lex.Text()
		if field.Name == "" {
			field.Name = token

			lex.Next()
			continue
		}
		// }

		// consume {
		lex.Next()
		// }
	}

	// consume ';' {
	lex.Next()
	// }

	return field
}
