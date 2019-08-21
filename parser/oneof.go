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
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (f *OneofField) SetInlineComment(comment *Comment) {
	f.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (f *OneofField) Accept(v Visitor) {
	if !v.VisitOneofField(f) {
		return
	}

	for _, comment := range f.Comments {
		comment.Accept(v)
	}
	if f.InlineComment != nil {
		f.InlineComment.Accept(v)
	}
}

// Oneof consists of oneof fields and a oneof name.
type Oneof struct {
	OneofFields []*OneofField
	OneofName   string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// InlineCommentBehindLeftCurly is the optional one placed behind a left curly.
	InlineCommentBehindLeftCurly *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (o *Oneof) SetInlineComment(comment *Comment) {
	o.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (o *Oneof) Accept(v Visitor) {
	if !v.VisitOneof(o) {
		return
	}

	for _, field := range o.OneofFields {
		field.Accept(v)
	}
	for _, comment := range o.Comments {
		comment.Accept(v)
	}
	if o.InlineComment != nil {
		o.InlineComment.Accept(v)
	}
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

	inlineLeftCurly := p.parseInlineComment()

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
		p.MaybeScanInlineComment(oneofField)
		oneofFields = append(oneofFields, oneofField)

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			break
		} else {
			p.lex.UnNext()
		}
	}

	if p.permissive {
		// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/issues/30.
		p.lex.ConsumeToken(scanner.TSEMICOLON)
	}

	return &Oneof{
		OneofFields:                  oneofFields,
		OneofName:                    oneofName,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta:                         meta.NewMetaWithLastPos(startPos, p.lex.Pos),
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
