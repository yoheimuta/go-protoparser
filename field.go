package protoparser

import "text/scanner"

// Field is the basic element of a protocol buffer message.
type Field struct {
	Comments []string
	Type     *Type
	Name     string
}

// type name = number validator';'
func parseField(lex *lexer) *Field {
	field := &Field{}

	for lex.text() != ";" && lex.token != scanner.EOF {
		token := lex.text()
		if field.Type == nil {
			field.Type = parseType(lex)
			continue
		}
		if field.Name == "" {
			field.Name = token

			lex.next()
			continue
		}
		// consume {
		lex.next()
		// }
	}

	// consume ';' {
	lex.next()
	// }

	return field
}
