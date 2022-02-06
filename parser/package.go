package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Package can be used to prevent name clashes between protocol message types.
type Package struct {
	Name string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (p *Package) SetInlineComment(comment *Comment) {
	p.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (p *Package) Accept(v Visitor) {
	if !v.VisitPackage(p) {
		return
	}

	for _, comment := range p.Comments {
		comment.Accept(v)
	}
	if p.InlineComment != nil {
		p.InlineComment.Accept(v)
	}
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
	startPos := p.lex.Pos

	ident, _, err := p.lex.ReadFullIdent()
	if err != nil {
		return nil, p.unexpected("fullIdent")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Package{
		Name: ident,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
