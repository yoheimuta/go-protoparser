package protoparser

import "fmt"

// Oneof consists of oneof fields and a oneof name.
type Oneof struct {
	Comments []string
	Name     string
	Fields   []*Field
}

// "oneof" var '{' OneofContent '}'
func parseOneof(lex *lexer) (*Oneof, error) {
	text := lex.text()
	if text != "oneof" {
		return nil, fmt.Errorf("not found oneof, text=%s", text)
	}

	// get the name {
	lex.next()
	name := lex.text()
	lex.next()
	// }

	// get the content {
	/// consume '{' {
	lex.next()
	/// }
	fields, _, _, _, err := parseMessageContent(lex)
	if err != nil {
		return nil, err
	}

	/// consume '}' {
	lex.next()
	/// }
	// }

	return &Oneof{
		Name:   name,
		Fields: fields,
	}, nil
}
