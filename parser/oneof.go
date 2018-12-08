package parser

import (
	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

// OneofField is a constituent field of oneof.
type OneofField struct {
	Type         string
	FieldName    string
	FieldNumber  string
	FieldOptions []*FieldOption

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// Oneof consists of oneof fields and a oneof name.
type Oneof struct {
	OneofFields []*OneofField
	OneofName   string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (o *Oneof) SetInlineComment(comment *Comment) {
	o.InlineComment = comment
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
	startPos := p.lex.Pos

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
		comments := p.ParseComments()

		err := p.lex.ReadEmptyStatement()
		if err == nil {
			continue
		}

		oneofField, err := p.parseOneofField()
		if err != nil {
			return nil, err
		}
		oneofField.Comments = comments
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
		Meta:        meta.NewMeta(startPos),
	}, nil
}

// oneofField = type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
// https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#oneof_and_oneof_field
func (p *Parser) parseOneofField() (*OneofField, error) {
	typeValue, startPos, err := p.parseType()
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
		Meta:         meta.NewMeta(startPos),
	}, nil
}
