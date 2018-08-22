package protoparser

import (
	"text/scanner"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

func parseComments(lex *lexer.Lexer) []string {
	var s []string
	for lex.Token == scanner.Comment {
		s = append(s, lex.Text())
		lex.Next()
	}
	return s
}
