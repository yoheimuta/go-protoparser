package protoparser

import (
	"io"
	"text/scanner"
)

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func newlexer(input io.Reader) *lexer {
	lex := new(lexer)
	lex.scan.Init(input)
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanComments
	lex.next()
	return lex
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }
