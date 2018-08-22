package protoparser

import "text/scanner"

func parseComments(lex *Lexer) []string {
	var s []string
	for lex.token == scanner.Comment {
		s = append(s, lex.Text())
		lex.Next()
	}
	return s
}
