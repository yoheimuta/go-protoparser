package protoparser

import (
	"strings"
	"text/scanner"
)

func lex(input string) *lexer {
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanComments
	lex.next()
	return lex
}
