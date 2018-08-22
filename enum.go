package protoparser

import "fmt"

// Enum consists of a name and an enum body.
type Enum struct {
	Comments   []string
	Name       string
	EnumFields []*EnumField
}

// "enum" var '{' EnumContent '}'
func parseEnum(lex *Lexer) (*Enum, error) {
	text := lex.Text()
	if text != "enum" {
		return nil, fmt.Errorf("not found enum, Text=%s", text)
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
	fields, err := parseEnumContent(lex)
	if err != nil {
		return nil, err
	}
	// }

	// consume '}' {
	lex.Next()
	// }

	return &Enum{
		Name:       name,
		EnumFields: fields,
	}, nil
}

// EnumField...}
func parseEnumContent(lex *Lexer) ([]*EnumField, error) {
	var fields []*EnumField

	for lex.Text() != "}" {
		field, err := parseEnumField(lex)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	return fields, nil
}
