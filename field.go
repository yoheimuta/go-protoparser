package protoparser

import "text/scanner"

// Field is the basic element of a protocol buffer message.
type Field struct {
	Comments    []string
	Type        *Type
	Name        string
	HasRepeated bool
}

// type name = number validator';'
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#normal_field
func parseField(lex *lexer) *Field {
	field := &Field{}

	// check HasRepeated {
	text := lex.text()
	if text == "repeated" {
		field.HasRepeated = true
		lex.next()
	}
	// }

	for lex.text() != ";" && lex.token != scanner.EOF {
		// get the type {
		if field.Type == nil {
			field.Type = parseType(lex)
			continue
		}
		// }

		// get the name {
		token := lex.text()
		if field.Name == "" {
			field.Name = token

			lex.next()
			continue
		}
		// }

		// consume {
		lex.next()
		// }
	}

	// consume ';' {
	lex.next()
	// }

	return field
}
