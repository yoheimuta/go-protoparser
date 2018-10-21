package parser

import (
	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

// OneofField is a constituent field of oneof.
type OneofField struct {
	Type         string
	FieldName    string
	FieldNumber  string
	FieldOptions []*FieldOption
}

// Oneof consists of oneof fields and a oneof name.
type Oneof struct {
	OneofFields []*OneofField
	OneofName   string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
}

// ParseOneof parses the oneof.
//  oneof = "oneof" oneofName "{" { oneofField | emptyStatement } "}"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#oneof_and_oneof_field
func (p *Parser) ParseOneof() (*Oneof, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TONEOF {
		return nil, p.unexpected("oneof")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("oneofName")
	}
	oneofName := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, p.unexpected("{")
	}

	var oneofFields []*OneofField
	for {
		err := p.lex.ReadEmptyStatement()
		if err == nil {
			continue
		}

		oneofField, err := p.parseOneofField()
		if err != nil {
			return nil, err
		}
		oneofFields = append(oneofFields, oneofField)

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			break
		} else {
			p.lex.UnNext()
		}
	}

	return &Oneof{
		OneofFields: oneofFields,
		OneofName:   oneofName,
	}, nil
}

// oneofField = type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
// https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#oneof_and_oneof_field
func (p *Parser) parseOneofField() (*OneofField, error) {
	typeValue, err := p.parseType()
	if err != nil {
		return nil, p.unexpected("type")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("fieldName")
	}
	fieldName := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	fieldNumber, err := p.parseFieldNumber()
	if err != nil {
		return nil, p.unexpected("fieldNumber")
	}

	fieldOptions, err := p.parseFieldOptionsOption()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &OneofField{
		Type:         typeValue,
		FieldName:    fieldName,
		FieldNumber:  fieldNumber,
		FieldOptions: fieldOptions,
	}, nil
}
