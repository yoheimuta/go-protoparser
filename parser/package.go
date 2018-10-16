package parser

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// ParsePackage parses the package.
// package = "package" fullIdent ";"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#package
func (p *Parser) ParsePackage() (string, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TPACKAGE {
		return "", p.unexpected("package")
	}

	ident, err := p.lex.ReadFullIdent()
	if err != nil {
		return "", p.unexpected("fullIdent")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return "", p.unexpected(";")
	}

	return ident, nil
}
