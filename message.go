package protoparser

import (
	"fmt"
	"text/scanner"
	"github.com/yoheimuta/go-protoparser/internal/lexer"
)

// Message consists of a message name and a message body.
type Message struct {
	Comments []string
	Name     string
	Fields   []*Field
	Messages []*Message
	Enums    []*Enum
	Oneofs   []*Oneof
}

// "message" var '{' messageContent '}'
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#message_definition
func parseMessage(lex *lexer.Lexer) (*Message, error) {
	text := lex.Text()
	if text != "message" {
		return nil, fmt.Errorf("not found message, Text=%s", text)
	}

	// get name {
	lex.Next()
	name := lex.Text()
	lex.Next()
	// }

	// get content {
	/// consume '{' {
	lex.Next()
	/// }
	fields, nests, enums, oneofs, err := parseMessageContent(lex)
	if err != nil {
		return nil, err
	}
	// }

	// consume '}' {
	lex.Next()
	// }

	return &Message{
		Name:     name,
		Fields:   fields,
		Messages: nests,
		Enums:    enums,
		Oneofs:   oneofs,
	}, nil
}

// "message" ...
// "enum" ...
// "oneof" ...
// field
func parseMessageContent(lex *lexer.Lexer) (fields []*Field, messages []*Message, enums []*Enum, oneofs []*Oneof, err error) {
	for lex.Text() != "}" {
		if lex.Token != scanner.Comment {
			return nil, nil, nil, nil, fmt.Errorf("not found comment, Text=%s", lex.Text())
		}
		comments := parseComments(lex)

		switch lex.Text() {
		case "message":
			message, parseErr := parseMessage(lex)
			if parseErr != nil {
				return nil, nil, nil, nil, parseErr
			}
			message.Comments = append(message.Comments, comments...)
			messages = append(messages, message)
		case "enum":
			enum, parseErr := parseEnum(lex)
			if parseErr != nil {
				return nil, nil, nil, nil, parseErr
			}
			enum.Comments = append(enum.Comments, comments...)
			enums = append(enums, enum)
		case "oneof":
			oneof, parseErr := parseOneof(lex)
			if parseErr != nil {
				return nil, nil, nil, nil, parseErr
			}
			oneof.Comments = append(oneof.Comments, comments...)
			oneofs = append(oneofs, oneof)
		default:
			field := parseField(lex)
			field.Comments = append(field.Comments, comments...)
			fields = append(fields, field)
		}
	}

	return fields, messages, enums, oneofs, nil
}
