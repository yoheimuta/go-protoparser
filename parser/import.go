package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// ImportModifier is a modifier enum type for import behavior.
type ImportModifier uint

// Optional import modifier value to change a default behavior.
const (
	ImportModifierNone ImportModifier = iota
	ImportModifierPublic
	ImportModifierWeak
)

// Import is used to import another .proto's definitions.
type Import struct {
	Modifier ImportModifier
	Location string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (i *Import) SetInlineComment(comment *Comment) {
	i.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (i *Import) Accept(v Visitor) {
	if !v.VisitImport(i) {
		return
	}

	for _, comment := range i.Comments {
		comment.Accept(v)
	}
	if i.InlineComment != nil {
		i.InlineComment.Accept(v)
	}
}

// ParseImport parses the import.
//  import = "import" [ "weak" | "public" ] strLit ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#import_statement
func (p *Parser) ParseImport() (*Import, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TIMPORT {
		return nil, p.unexpected(`"import"`)
	}
	startPos := p.lex.Pos

	var modifier ImportModifier
	p.lex.NextKeywordOrStrLit()
	switch p.lex.Token {
	case scanner.TPUBLIC:
		modifier = ImportModifierPublic
	case scanner.TWEAK:
		modifier = ImportModifierWeak
	case scanner.TSTRLIT:
		modifier = ImportModifierNone
		p.lex.UnNext()
	}

	p.lex.NextStrLit()
	if p.lex.Token != scanner.TSTRLIT {
		return nil, p.unexpected("strLit")
	}
	location := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Import{
		Modifier: modifier,
		Location: location,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
