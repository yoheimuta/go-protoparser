package protoparser

import (
	"fmt"
	"text/scanner"
)

// Message consists of a message name and a message body.
type Message struct {
	Comments []string
	Name     string
	Fields   []*Field
	Nests    []*Message
	Enums    []*Enum
	Oneofs   []*Oneof
}

// "message" var '{' messageContent '}'
func parseMessage(lex *lexer) (*Message, error) {
	text := lex.text()
	if text != "message" {
		return nil, fmt.Errorf("not found message, text=%s", text)
	}

	// get name {
	lex.next()
	name := lex.text()
	lex.next()
	// }

	// get content {
	/// consume '{' {
	lex.next()
	/// }
	fields, nests, enums, oneofs, err := parseMessageContent(lex)
	if err != nil {
		return nil, err
	}
	// }

	// consume '}' {
	lex.next()
	// }

	return &Message{
		Name:   name,
		Fields: fields,
		Nests:  nests,
		Enums:  enums,
		Oneofs: oneofs,
	}, nil
}

// "message"
// "enum"
// field
func parseMessageContent(lex *lexer) (fields []*Field, messages []*Message, enums []*Enum, oneofs []*Oneof, err error) {
	for lex.text() != "}" {
		if lex.token != scanner.Comment {
			return nil, nil, nil, nil, fmt.Errorf("not found comment, text=%s", lex.text())
		}
		comments := parseComments(lex)

		switch lex.text() {
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
