package parser

import (
	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

// Option can be used in proto files, messages, enums and services.
type Option struct {
	OptionName string
	Constant   string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// ParseOption parses the option.
//  option = "option" optionName  "=" constant ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#option
func (p *Parser) ParseOption() (*Option, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TOPTION {
		return nil, p.unexpected("option")
	}
	startPos := p.lex.Pos

	optionName, err := p.parseOptionName()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	constant, _, err := p.lex.ReadConstant()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Option{
		OptionName: optionName,
		Constant:   constant,
		Meta:       meta.NewMeta(startPos),
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
		fullIdent, _, err := p.lex.ReadFullIdent()
		if err != nil {
			return "", err
		}
		optionName += fullIdent

		p.lex.Next()
		if p.lex.Token != scanner.TRIGHTPAREN {
			return "", p.unexpected(")")
		}
		optionName += p.lex.Text
	default:
		return "", p.unexpected("ident or left paren")
	}

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TDOT {
			p.lex.UnNext()
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
