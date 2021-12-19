package parser

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

type parseMessageBodyStatementErr struct {
	parseFieldErr          error
	parseEmptyStatementErr error
}

func (e *parseMessageBodyStatementErr) Error() string {
	return fmt.Sprintf(
		"%v:%v",
		e.parseFieldErr,
		e.parseEmptyStatementErr,
	)
}

// Message consists of a message name and a message body.
type Message struct {
	MessageName string
	// MessageBody can have fields, nested enum definitions, nested message definitions,
	// options, oneofs, map fields, group fields(proto2 only), extends, reserved, and extensions(proto2 only) statements.
	MessageBody []Visitee

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
func (m *Message) SetInlineComment(comment *Comment) {
	m.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (m *Message) Accept(v Visitor) {
	if !v.VisitMessage(m) {
		return
	}

	for _, body := range m.MessageBody {
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

// ParseMessage parses the message.
//  message = "message" messageName messageBody
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#message_definition
func (p *Parser) ParseMessage() (*Message, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TMESSAGE {
		return nil, p.unexpected("message")
	}
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("messageName")
	}
	messageName := p.lex.Text

	messageBody, inlineLeftCurly, lastPos, err := p.parseMessageBody()
	if err != nil {
		return nil, err
	}

	return &Message{
		MessageName:                  messageName,
		MessageBody:                  messageBody,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

// messageBody = "{" { field | enum | message | option | oneof | mapField | reserved | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#message_definition
func (p *Parser) parseMessageBody() (
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
		case scanner.TENUM:
			enum, err := p.ParseEnum()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			enum.Comments = comments
			stmt = enum
		case scanner.TMESSAGE:
			message, err := p.ParseMessage()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			message.Comments = comments
			stmt = message
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TONEOF:
			oneof, err := p.ParseOneof()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			oneof.Comments = comments
			stmt = oneof
		case scanner.TMAP:
			mapField, err := p.ParseMapField()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			mapField.Comments = comments
			stmt = mapField
		case scanner.TEXTEND:
			extend, err := p.ParseExtend()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			extend.Comments = comments
			stmt = extend
		case scanner.TRESERVED:
			reserved, err := p.ParseReserved()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			reserved.Comments = comments
			stmt = reserved
		case scanner.TEXTENSIONS:
			extensions, err := p.ParseExtensions()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			extensions.Comments = comments
			stmt = extensions
		default:
			var ferr error
			isGroup := p.peekIsGroup()
			if isGroup {
				groupField, groupErr := p.ParseGroupField()
				if groupErr == nil {
					groupField.Comments = comments
					stmt = groupField
					break
				}
				ferr = groupErr
				p.lex.UnNext()
			} else {
				field, fieldErr := p.ParseField()
				if fieldErr == nil {
					field.Comments = comments
					stmt = field
					break
				}
				ferr = fieldErr
				p.lex.UnNext()
			}

			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmt = &EmptyStatement{}
				break
			}

			return nil, nil, scanner.Position{}, &parseMessageBodyStatementErr{
				parseFieldErr:          ferr,
				parseEmptyStatementErr: emptyErr,
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
	}
}
