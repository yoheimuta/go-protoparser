package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Edition is used to define the protobuf version.
type Edition struct {
	Edition string

	// EditionQuote includes quotes
	EditionQuote string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (s *Edition) SetInlineComment(comment *Comment) {
	s.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (s *Edition) Accept(v Visitor) {
	if !v.VisitEdition(s) {
		return
	}

	for _, comment := range s.Comments {
		comment.Accept(v)
	}
	if s.InlineComment != nil {
		s.InlineComment.Accept(v)
	}
}

// ParseEdition parses the Edition.
//
// edition = "edition" "=" [ ( "'" decimalLit "'" ) | ( '"' decimalLit '"' ) ] ";"
//
// See https://protobuf.dev/reference/protobuf/edition-2023-spec/#edition
func (p *Parser) ParseEdition() (*Edition, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TEDITION {
		p.lex.UnNext()
		return nil, nil
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

	p.lex.NextNumberLit()
	if p.lex.Token != scanner.TINTLIT {
		return nil, p.unexpected("intLit")
	}
	edition := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TQUOTE {
		return nil, p.unexpected("quote")
	}
	tq := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Edition{
		Edition:      edition,
		EditionQuote: lq + edition + tq,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
