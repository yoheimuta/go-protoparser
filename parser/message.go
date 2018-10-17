package parser

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

// ParseMessageBodyStatementErr consists of the enum body statement parse errors.
type ParseMessageBodyStatementErr struct {
	ParseFieldErr          error
	ParseEmptyStatementErr error
}

// Error represents an error condition.
func (e *ParseMessageBodyStatementErr) Error() string {
	return fmt.Sprintf(
		"%v:%v",
		e.ParseFieldErr,
		e.ParseEmptyStatementErr,
	)
}

// Message consists of a message name and a message body.
type Message struct {
	MessageName string
	// MessageBody can have fields, nested enum definitions, nested message definitions,
	// options, oneofs, map fields, and reserved statements.
	MessageBody []interface{}
}

// ParseMessage parses the message.
// message = "message" messageName messageBody
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#message_definition
func (p *Parser) ParseMessage() (*Message, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TMESSAGE {
		return nil, p.unexpected("message")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("messageName")
	}
	messageName := p.lex.Text

	messageBody, err := p.parseMessageBody()
	if err != nil {
		return nil, err
	}

	return &Message{
		MessageName: messageName,
		MessageBody: messageBody,
	}, nil
}

// messageBody = "{" { field | enum | message | option | oneof | mapField | reserved | emptyStatement } "}"
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#message_definition
func (p *Parser) parseMessageBody() ([]interface{}, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, p.unexpected("{")
	}

	var stmts []interface{}

	for {
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		switch token {
		case scanner.TENUM:
			enum, err := p.ParseEnum()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, enum)
		case scanner.TMESSAGE:
			message, err := p.ParseMessage()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, message)
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, option)
		case scanner.TONEOF:
			oneof, err := p.ParseOneof()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, oneof)
		case scanner.TMAP:
			mapField, err := p.ParseMapField()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, mapField)
		case scanner.TRESERVED:
			reserved, err := p.ParseReserved()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, reserved)
		default:
			field, fieldErr := p.ParseField()
			if fieldErr == nil {
				stmts = append(stmts, field)
				break
			}
			p.lex.UnNext()

			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmts = append(stmts, EmptyStatement{})
				break
			}

			return nil, &ParseMessageBodyStatementErr{
				ParseFieldErr:          fieldErr,
				ParseEmptyStatementErr: emptyErr,
			}
		}

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return stmts, nil
		}
		p.lex.UnNext()
	}
}
