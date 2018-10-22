package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// Package can be used to prevent name clashes between protocol message types.
type Package struct {
	Name string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
}

// ParsePackage parses the package.
//  package = "package" fullIdent ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#package
func (p *Parser) ParsePackage() (*Package, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TPACKAGE {
		return nil, p.unexpected("package")
	}

	ident, err := p.lex.ReadFullIdent()
	if err != nil {
		return nil, p.unexpected("fullIdent")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Package{
		Name: ident,
	}, nil
}
