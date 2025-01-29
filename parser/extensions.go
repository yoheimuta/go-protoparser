package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Extensions declare that a range of field numbers in a message are available for third-party extensions.
type Extensions struct {
	Ranges       []*Range
	Declarations []*Declaration

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// InlineCommentBehindLeftSquare is the optional one placed behind a left square.
	InlineCommentBehindLeftSquare *Comment
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

	for _, declaration := range e.Declarations {
		declaration.Accept(v)
	}
	for _, comment := range e.Comments {
		comment.Accept(v)
	}
	if e.InlineComment != nil {
		e.InlineComment.Accept(v)
	}
}

// ParseExtensions parses the extensions.
//
//	extensions = "extensions" ranges ";"
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

	declarations, inlineLeftSquare, err := p.parseDeclarations()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Extensions{
		Ranges:                        ranges,
		Declarations:                  declarations,
		InlineCommentBehindLeftSquare: inlineLeftSquare,
		Meta:                          meta.Meta{Pos: startPos.Position, LastPos: p.lex.Pos.Position},
	}, nil
}

// parseDeclarations parses the declarations.
//
//	declarations = "[" declaration { ","  declaration } "]"
//
// See https://protobuf.dev/programming-guides/extension_declarations/
func (p *Parser) parseDeclarations() ([]*Declaration, *Comment, error) {
	declarations := []*Declaration{}
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTSQUARE {
		p.lex.UnNext()
		return nil, nil, nil
	}
	inlineLeftSquare := p.parseInlineComment()

	for {
		comments := p.ParseComments()

		declaration, err := p.ParseDeclaration()
		if err != nil {
			return nil, nil, err
		}
		declaration.Comments = comments
		declarations = append(declarations, declaration)

		p.lex.Next()
		token := p.lex.Token
		inlineComment1 := p.parseInlineComment()
		if token == scanner.TRIGHTSQUARE {
			p.assignInlineComments(declaration, inlineComment1, p.parseInlineComment())
			break
		}
		if token != scanner.TCOMMA {
			return nil, nil, p.unexpected(", or ]")
		}
		p.assignInlineComments(declaration, inlineComment1, p.parseInlineComment())
	}
	return declarations, inlineLeftSquare, nil
}

// assignInlineComments assigns inline comments to a declaration, ensuring proper order.
func (p *Parser) assignInlineComments(declaration *Declaration, comments ...*Comment) {
	for _, comment := range comments {
		if comment != nil {
			declaration.SetInlineComment(comment)
		}
	}
}
