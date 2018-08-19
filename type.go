package protoparser

// Type is a field type.
type Type struct {
	// Name is a type name.
	Name string
}

func parseType(lex *lexer) *Type {
	s := lex.text()
	lex.next()

	for lex.text() == "." {
		s += lex.text()
		lex.next()
		s += lex.text()
		lex.next()
	}
	return &Type{
		Name: s,
	}
}
