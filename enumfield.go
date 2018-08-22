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
func parseEnumField(lex *Lexer) (*EnumField, error) {
	field := &EnumField{}

	// get comments {
	if lex.token != scanner.Comment {
		return nil, fmt.Errorf("not found comment, Text=%s", lex.Text())
	}
	field.Comments = parseComments(lex)
	// }

	field.Name = lex.Text()

	// consume the rest {
	for lex.Text() != ";" && lex.token != scanner.EOF {
		lex.Next()
	}
	lex.Next()
	// }
	return field, nil
}
