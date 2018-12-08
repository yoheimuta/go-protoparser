package parser

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/parser/meta"
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
	// options, oneofs, map fields, and reserved statements.
	MessageBody []interface{}

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

	messageBody, inlineLeftCurly, err := p.parseMessageBody()
	if err != nil {
		return nil, err
	}

	return &Message{
		MessageName:                  messageName,
		MessageBody:                  messageBody,
		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta:                         meta.NewMeta(startPos),
	}, nil
}

// messageBody = "{" { field | enum | message | option | oneof | mapField | reserved | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#message_definition
func (p *Parser) parseMessageBody() ([]interface{}, *Comment, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	// Parses emptyBody. This spec is not documented, but allowed in general. {
	p.lex.Next()
	if p.lex.Token == scanner.TRIGHTCURLY {
		return nil, nil, nil
	}
	p.lex.UnNext()
	// }

	var stmts []interface{}

	for {
		comments := p.ParseComments()

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		var stmt interface {
			HasInlineCommentSetter
		}

		switch token {
		case scanner.TENUM:
			enum, err := p.ParseEnum()
			if err != nil {
				return nil, nil, err
			}
			enum.Comments = comments
			stmt = enum
		case scanner.TMESSAGE:
			message, err := p.ParseMessage()
			if err != nil {
				return nil, nil, err
			}
			message.Comments = comments
			stmt = message
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TONEOF:
			oneof, err := p.ParseOneof()
			if err != nil {
				return nil, nil, err
			}
			oneof.Comments = comments
			stmt = oneof
		case scanner.TMAP:
			mapField, err := p.ParseMapField()
			if err != nil {
				return nil, nil, err
			}
			mapField.Comments = comments
			stmt = mapField
		case scanner.TRESERVED:
			reserved, err := p.ParseReserved()
			if err != nil {
				return nil, nil, err
			}
			reserved.Comments = comments
			stmt = reserved
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

			return nil, nil, &parseMessageBodyStatementErr{
				parseFieldErr:          fieldErr,
				parseEmptyStatementErr: emptyErr,
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return stmts, inlineLeftCurly, nil
		}
		p.lex.UnNext()
	}
}
