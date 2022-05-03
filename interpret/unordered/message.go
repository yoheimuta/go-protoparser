package unordered

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// MessageBody is unordered in nature, but each slice field preserves the original order.
type MessageBody struct {
	Fields          []*parser.Field
	Enums           []*Enum
	Messages        []*Message
	Options         []*parser.Option
	Oneofs          []*parser.Oneof
	Maps            []*parser.MapField
	Groups          []*parser.GroupField
	Reserves        []*parser.Reserved
	Extends         []*parser.Extend
	EmptyStatements []*parser.EmptyStatement
}

// Message consists of a message name and a message body.
type Message struct {
	MessageName string
	MessageBody *MessageBody

	// Comments are the optional ones placed at the beginning.
	Comments []*parser.Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *parser.Comment
	// InlineCommentBehindLeftCurly is the optional one placed behind a left curly.
	InlineCommentBehindLeftCurly *parser.Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// InterpretMessage interprets *parser.Message to *Message.
func InterpretMessage(src *parser.Message) (*Message, error) {
	if src == nil {
		return nil, nil
	}

	messageBody, err := interpretMessageBody(src.MessageBody)
	if err != nil {
		return nil, fmt.Errorf("invalid Message %s: %w", src.MessageName, err)
	}
	return &Message{
		MessageName:                  src.MessageName,
		MessageBody:                  messageBody,
		Comments:                     src.Comments,
		InlineComment:                src.InlineComment,
		InlineCommentBehindLeftCurly: src.InlineCommentBehindLeftCurly,
		Meta:                         src.Meta,
	}, nil
}

func interpretMessageBody(src []parser.Visitee) (
	*MessageBody,
	error,
) {
	var fields []*parser.Field
	var enums []*Enum
	var messages []*Message
	var options []*parser.Option
	var oneofs []*parser.Oneof
	var maps []*parser.MapField
	var groups []*parser.GroupField
	var reserves []*parser.Reserved
	var extends []*parser.Extend
	var emptyStatements []*parser.EmptyStatement
	for _, s := range src {
		switch t := s.(type) {
		case *parser.Field:
			fields = append(fields, t)
		case *parser.Enum:
			enum, err := InterpretEnum(t)
			if err != nil {
				return nil, err
			}
			enums = append(enums, enum)
		case *parser.Message:
			message, err := InterpretMessage(t)
			if err != nil {
				return nil, err
			}
			messages = append(messages, message)
		case *parser.Option:
			options = append(options, t)
		case *parser.Oneof:
			oneofs = append(oneofs, t)
		case *parser.MapField:
			maps = append(maps, t)
		case *parser.GroupField:
			groups = append(groups, t)
		case *parser.Reserved:
			reserves = append(reserves, t)
		case *parser.Extend:
			extends = append(extends, t)
		case *parser.EmptyStatement:
			emptyStatements = append(emptyStatements, t)
		default:
			return nil, fmt.Errorf("invalid MessageBody type %T of %v", t, t)
		}
	}
	return &MessageBody{
		Fields:          fields,
		Enums:           enums,
		Messages:        messages,
		Options:         options,
		Oneofs:          oneofs,
		Maps:            maps,
		Groups:          groups,
		Reserves:        reserves,
		Extends:         extends,
		EmptyStatements: emptyStatements,
	}, nil
}
