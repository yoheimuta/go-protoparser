package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Extensions declare that a range of field numbers in a message are available for third-party extensions.
type Extensions struct {
	Ranges []*Range

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (e *Extensions) SetInlineComment(comment *Comment) {
	e.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (e *Extensions) Accept(v Visitor) {
	if !v.VisitExtensions(e) {
		return
	}

	for _, comment := range e.Comments {
		comment.Accept(v)
	}
	if e.InlineComment != nil {
		e.InlineComment.Accept(v)
	}
}

// ParseExtensions parses the extensions.
//  extensions = "extensions" ranges ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto2-spec#extensions
func (p *Parser) ParseExtensions() (*Extensions, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TEXTENSIONS {
		return nil, p.unexpected("extensions")
	}
	startPos := p.lex.Pos

	ranges, err := p.parseRanges()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Extensions{
		Ranges: ranges,
		Meta:   meta.Meta{Pos: startPos.Position},
	}, nil
}
