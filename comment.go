package protoparser

import "text/scanner"

func parseComments(lex *lexer) []string {
	var s []string
	for lex.token == scanner.Comment {
		s = append(s, lex.text())
		lex.next()
	}
	return s
}
