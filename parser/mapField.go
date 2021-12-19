package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// MapField is an associative map.
type MapField struct {
	KeyType      string
	Type         string
	MapName      string
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
func (m *MapField) SetInlineComment(comment *Comment) {
	m.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (m *MapField) Accept(v Visitor) {
	if !v.VisitMapField(m) {
		return
	}

	for _, comment := range m.Comments {
		comment.Accept(v)
	}
	if m.InlineComment != nil {
		m.InlineComment.Accept(v)
	}
}

// ParseMapField parses the mapField.
//  mapField = "map" "<" keyType "," type ">" mapName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#map_field
func (p *Parser) ParseMapField() (*MapField, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TMAP {
		return nil, p.unexpected("map")
	}
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TLESS {
		return nil, p.unexpected("<")
	}

	keyType, err := p.parseKeyType()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TCOMMA {
		return nil, p.unexpected(",")
	}

	typeValue, _, err := p.parseType()
	if err != nil {
		return nil, p.unexpected("type")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TGREATER {
		return nil, p.unexpected(">")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("mapName")
	}
	mapName := p.lex.Text

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

	return &MapField{
		KeyType:      keyType,
		Type:         typeValue,
		MapName:      mapName,
		FieldNumber:  fieldNumber,
		FieldOptions: fieldOptions,
		Meta:         meta.Meta{Pos: startPos.Position},
	}, nil
}

var keyTypeConstants = map[string]struct{}{
	"int32":    {},
	"int64":    {},
	"uint32":   {},
	"uint64":   {},
	"sint32":   {},
	"sint64":   {},
	"fixed32":  {},
	"fixed64":  {},
	"sfixed32": {},
	"sfixed64": {},
	"bool":     {},
	"string":   {},
}

// keyType = "int32" | "int64" | "uint32" | "uint64" | "sint32" | "sint64" |
//          "fixed32" | "fixed64" | "sfixed32" | "sfixed64" | "bool" | "string"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#map_field
func (p *Parser) parseKeyType() (string, error) {
	p.lex.Next()
	if _, ok := keyTypeConstants[p.lex.Text]; ok {
		return p.lex.Text, nil
	}
	return "", p.unexpected("keyType constant")
}
