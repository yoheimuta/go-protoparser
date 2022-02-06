package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Syntax is used to define the protobuf version.
type Syntax struct {
	ProtobufVersion string

	// ProtobufVersionQuote includes quotes
	ProtobufVersionQuote string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (s *Syntax) SetInlineComment(comment *Comment) {
	s.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (s *Syntax) Accept(v Visitor) {
	if !v.VisitSyntax(s) {
		return
	}

	for _, comment := range s.Comments {
		comment.Accept(v)
	}
	if s.InlineComment != nil {
		s.InlineComment.Accept(v)
	}
}

// Version returns the version number.
func (s *Syntax) Version() int {
	switch s.ProtobufVersion {
	case "proto3":
		return 3
	case "proto2":
		return 2
	default:
		return 0
	}
}

// ParseSyntax parses the syntax.
//  syntax = "syntax" "=" quote "proto3" quote ";"
//  syntax = "syntax" "=" quote "proto2" quote ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#syntax
func (p *Parser) ParseSyntax() (*Syntax, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TSYNTAX {
		return nil, p.unexpected("syntax")
	}
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TQUOTE {
		return nil, p.unexpected("quote")
	}
	lq := p.lex.Text

	p.lex.Next()
	if p.lex.Text != "proto3" && p.lex.Text != "proto2" {
		return nil, p.unexpected("proto3 or proto2")
	}
	version := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TQUOTE {
		return nil, p.unexpected("quote")
	}
	tq := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Syntax{
		ProtobufVersion:      version,
		ProtobufVersionQuote: lq + version + tq,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
