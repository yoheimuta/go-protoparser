package parser

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Declaration is an option of extension ranges.
type Declaration struct {
	Number   string
	FullName string
	Type     string
	Reserved bool
	Repeated bool

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
func (d *Declaration) SetInlineComment(comment *Comment) {
	d.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (d *Declaration) Accept(v Visitor) {
	if !v.VisitDeclaration(d) {
		return
	}

	for _, comment := range d.Comments {
		comment.Accept(v)
	}
	if d.InlineComment != nil {
		d.InlineComment.Accept(v)
	}
	if d.InlineCommentBehindLeftCurly != nil {
		d.InlineCommentBehindLeftCurly.Accept(v)
	}
}

// ParseDeclaration parses a declaration.
//
//	declaration = "declaration" "=" "{"
//	  "number" ":" number ","
//	  "full_name" ":" string ","
//	  "type" ":" string ","
//	  "repeated" ":" bool ","
//	  "reserved" ":" bool
//	"}"
//
// See https://protobuf.dev/programming-guides/extension_declarations/
func (p *Parser) ParseDeclaration() (*Declaration, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TDECLARATION {
		return nil, p.unexpected("declaration")
	}
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	var number string
	var fullName string
	var typeStr string
	var repeated bool
	var reserved bool

	for {
		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			break
		}
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
		}

		p.lex.NextKeyword()
		if p.lex.Token == scanner.TNUMBER {
			p.lex.Next()
			if p.lex.Token != scanner.TCOLON {
				return nil, p.unexpected(":")
			}
			p.lex.NextNumberLit()
			if p.lex.Token != scanner.TINTLIT {
				return nil, p.unexpected("number")
			}
			number = p.lex.Text
		} else if p.lex.Token == scanner.TFULLNAME {
			p.lex.Next()
			if p.lex.Token != scanner.TCOLON {
				return nil, p.unexpected(":")
			}
			p.lex.NextStrLit()
			if p.lex.Token != scanner.TSTRLIT {
				return nil, p.unexpected("full_name string")
			}
			fullName = p.lex.Text
		} else if p.lex.Token == scanner.TTYPE {
			p.lex.Next()
			if p.lex.Token != scanner.TCOLON {
				return nil, p.unexpected(":")
			}
			p.lex.NextStrLit()
			if p.lex.Token != scanner.TSTRLIT {
				return nil, p.unexpected("type string")
			}
			typeStr = p.lex.Text
		} else if p.lex.Token == scanner.TREPEATED {
			p.lex.Next()
			if p.lex.Token != scanner.TCOLON {
				return nil, p.unexpected(":")
			}
			p.lex.Next()
			if p.lex.Token != scanner.TIDENT {
				return nil, p.unexpected("repeated bool")
			}
			repeated = p.lex.Text == "true"
		} else if p.lex.Token == scanner.TRESERVED {
			p.lex.Next()
			if p.lex.Token != scanner.TCOLON {
				return nil, p.unexpected(":")
			}
			p.lex.Next()
			if p.lex.Token != scanner.TIDENT {
				return nil, p.unexpected("reserved bool")
			}
			reserved = p.lex.Text == "true"
		} else {
			return nil, p.unexpected("number, full_name, type, repeated, reserved, or }")
		}
	}

	return &Declaration{
		Number:                       number,
		FullName:                     fullName,
		Type:                         typeStr,
		Reserved:                     reserved,
		Repeated:                     repeated,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta:                         meta.Meta{Pos: startPos.Position, LastPos: p.lex.Pos.Position},
	}, nil
}
