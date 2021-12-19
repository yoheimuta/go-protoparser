package parser

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

type parseExtendBodyStatementErr struct {
	parseFieldErr          error
	parseEmptyStatementErr error
}

func (e *parseExtendBodyStatementErr) Error() string {
	return fmt.Sprintf(
		"%v:%v",
		e.parseFieldErr,
		e.parseEmptyStatementErr,
	)
}

// Extend consists of a messageType and an extend body.
type Extend struct {
	MessageType string
	// ExtendBody can have fields and emptyStatements
	ExtendBody []Visitee

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
func (m *Extend) SetInlineComment(comment *Comment) {
	m.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (m *Extend) Accept(v Visitor) {
	if !v.VisitExtend(m) {
		return
	}

	for _, body := range m.ExtendBody {
		body.Accept(v)
	}
	for _, comment := range m.Comments {
		comment.Accept(v)
	}
	if m.InlineComment != nil {
		m.InlineComment.Accept(v)
	}
	if m.InlineCommentBehindLeftCurly != nil {
		m.InlineCommentBehindLeftCurly.Accept(v)
	}
}

// ParseExtend parses the extend.
// Note that group is not supported.
//  extend = "extend" messageType "{" {field | group | emptyStatement} "}"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto2-spec#extend
func (p *Parser) ParseExtend() (*Extend, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TEXTEND {
		return nil, p.unexpected("extend")
	}
	startPos := p.lex.Pos

	messageType, _, err := p.lex.ReadMessageType()
	if err != nil {
		return nil, err
	}

	extendBody, inlineLeftCurly, lastPos, err := p.parseExtendBody()
	if err != nil {
		return nil, err
	}

	return &Extend{
		MessageType:                  messageType,
		ExtendBody:                   extendBody,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

// extendBody = "{" {field | group | emptyStatement} "}"
func (p *Parser) parseExtendBody() (
	[]Visitee,
	*Comment,
	scanner.Position,
	error,
) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, scanner.Position{}, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	// Parses emptyBody. This spec is not documented, but allowed in general. {
	p.lex.Next()
	if p.lex.Token == scanner.TRIGHTCURLY {
		lastPos := p.lex.Pos
		if p.permissive {
			// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
			p.lex.ConsumeToken(scanner.TSEMICOLON)
			if p.lex.Token == scanner.TSEMICOLON {
				lastPos = p.lex.Pos
			}
		}

		return nil, nil, lastPos, nil
	}
	p.lex.UnNext()
	// }

	var stmts []Visitee

	for {
		comments := p.ParseComments()

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		var stmt interface {
			HasInlineCommentSetter
			Visitee
		}

		switch token {
		case scanner.TRIGHTCURLY:
			if p.bodyIncludingComments {
				for _, comment := range comments {
					stmts = append(stmts, Visitee(comment))
				}
			}
			p.lex.Next()

			lastPos := p.lex.Pos
			if p.permissive {
				// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
				p.lex.ConsumeToken(scanner.TSEMICOLON)
				if p.lex.Token == scanner.TSEMICOLON {
					lastPos = p.lex.Pos
				}
			}
			return stmts, inlineLeftCurly, lastPos, nil
		default:
			field, fieldErr := p.ParseField()
			if fieldErr == nil {
				field.Comments = comments
				stmt = field
				break
			}
			p.lex.UnNext()

			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmt = &EmptyStatement{}
				break
			}

			return nil, nil, scanner.Position{}, &parseExtendBodyStatementErr{
				parseFieldErr:          fieldErr,
				parseEmptyStatementErr: emptyErr,
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
	}
}
