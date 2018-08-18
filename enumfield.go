package protoparser

import (
	"fmt"
	"text/scanner"
)

// EnumField is the basic element of a protocol buffer enum..
type EnumField struct {
	Comments []string
	Name     string
}

// comment var '=' tag';'
func parseEnumField(lex *lexer) (*EnumField, error) {
	field := &EnumField{}

	// get comments {
	if lex.token != scanner.Comment {
		return nil, fmt.Errorf("not found comment, text=%s", lex.text())
	}
	field.Comments = parseComments(lex)
	// }

	field.Name = lex.text()

	// consume the rest {
	for lex.text() != ";" && lex.token != scanner.EOF {
		lex.next()
	}
	lex.next()
	// }
	return field, nil
}
