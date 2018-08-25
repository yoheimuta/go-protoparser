package parser

import (
	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

// Option can be used in proto files, messages, enums and services.
type Option struct {
	Name     string
	Constant string
}

// ParseOption parses the option.
// option = "option" optionName  "=" constant ";"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#option
func (p *Parser) ParseOption() (*Option, error) {
	p.lex.Next()
	if p.lex.Text != "option" {
		return nil, p.unexpected("option")
	}

	optionName, err := p.parseOptionName()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	constant, err := p.lex.ReadConstant()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Option{
		Name:     optionName,
		Constant: constant,
	}, nil
}

// optionName = ( ident | "(" fullIdent ")" ) { "." ident }
func (p *Parser) parseOptionName() (string, error) {
	var optionName string

	p.lex.Next()
	switch p.lex.Token {
	case scanner.TIDENT:
		optionName = p.lex.Text
	case scanner.TLEFTPAREN:
		optionName = p.lex.Text
		fullIdent, err := p.lex.ReadFullIdent()
		if err != nil {
			return "", err
		}
		optionName += fullIdent

		p.lex.Next()
		if p.lex.Token != scanner.TRIGHTPAREN {
			return "", p.unexpected(")")
		}
		optionName += p.lex.Text
	}

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TDOT {
			p.lex.SetIgnoreNext()
			break
		}
		optionName += p.lex.Text

		p.lex.Next()
		if p.lex.Token != scanner.TIDENT {
			return "", p.unexpected("ident")
		}
		optionName += p.lex.Text
	}
	return optionName, nil
}
