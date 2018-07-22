package protoparser

// Type is a field type.
type Type struct {
	// Name is a type name.
	Name string
	// IsRepeated represents whether a type with repeated or not.
	IsRepeated bool
}

func parseType(lex *lexer) *Type {
	s := lex.text()
	lex.next()
	if s == "repeated" {
		t := parseType(lex)
		return &Type{
			Name:       t.Name,
			IsRepeated: true,
		}
	}
	for lex.text() == "." {
		s += lex.text()
		lex.next()
		s += lex.text()
		lex.next()
	}
	return &Type{
		Name:       s,
		IsRepeated: false,
	}
}
