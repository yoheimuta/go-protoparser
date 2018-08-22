package protoparser

import "github.com/yoheimuta/go-protoparser/internal/lexer"

// Type is a field type.
type Type struct {
	// Name is a type name.
	Name string
}

func parseType(lex *lexer.Lexer) *Type {
	s := lex.Text()
	lex.Next()

	for lex.Text() == "." {
		s += lex.Text()
		lex.Next()
		s += lex.Text()
		lex.Next()
	}
	return &Type{
		Name: s,
	}
}
