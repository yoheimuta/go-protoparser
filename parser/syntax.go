package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// ParseSyntax parses the syntax.
// syntax = "syntax" "=" quote "proto3" quote ";"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#syntax
func (p *Parser) ParseSyntax() (string, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TSYNTAX {
		return "", p.unexpected("syntax")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return "", p.unexpected("=")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TQUOTE {
		return "", p.unexpected("quote")
	}

	p.lex.Next()
	if p.lex.Text != "proto3" {
		return "", p.unexpected("proto3")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TQUOTE {
		return "", p.unexpected("quote")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return "", p.unexpected(";")
	}

	return "proto3", nil
}
