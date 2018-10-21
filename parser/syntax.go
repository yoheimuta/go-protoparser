package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// Syntax is used to define the protobuf version.
type Syntax struct {
	ProtobufVersion string
	// Comments are the only ones placed at the beginning.
	Comments []*Comment
}

// ParseSyntax parses the syntax.
//  syntax = "syntax" "=" quote "proto3" quote ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#syntax
func (p *Parser) ParseSyntax() (*Syntax, error) {
	comments, _ := p.ParseComments()

	p.lex.NextKeyword()
	if p.lex.Token != scanner.TSYNTAX {
		return nil, p.unexpected("syntax")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TQUOTE {
		return nil, p.unexpected("quote")
	}

	p.lex.Next()
	if p.lex.Text != "proto3" {
		return nil, p.unexpected("proto3")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TQUOTE {
		return nil, p.unexpected("quote")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Syntax{
		ProtobufVersion: "proto3",
		Comments:        comments,
	}, nil
}
