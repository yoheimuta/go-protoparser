package protoparser

import "fmt"

// Enum consists of a name and an enum body.
type Enum struct {
	Comments   []string
	Name       string
	EnumFields []*EnumField
}

// "enum" var '{' EnumContent '}'
func parseEnum(lex *lexer) (*Enum, error) {
	text := lex.text()
	if text != "enum" {
		return nil, fmt.Errorf("not found enum, text=%s", text)
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
	fields, err := parseEnumContent(lex)
	if err != nil {
		return nil, err
	}
	// }

	// consume '}' {
	lex.next()
	// }

	return &Enum{
		Name:       name,
		EnumFields: fields,
	}, nil
}

// EnumField...}
func parseEnumContent(lex *lexer) ([]*EnumField, error) {
	var fields []*EnumField

	for lex.text() != "}" {
		field, err := parseEnumField(lex)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	return fields, nil
}
