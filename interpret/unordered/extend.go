package unordered

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// ExtendBody is unordered in nature, but each slice field preserves the original order.
type ExtendBody struct {
	Fields          []*parser.Field
	EmptyStatements []*parser.EmptyStatement
}

// Extend consists of a messageType and a extend body.
type Extend struct {
	MessageType string
	ExtendBody  *ExtendBody

	// Comments are the optional ones placed at the beginning.
	Comments []*parser.Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *parser.Comment
	// InlineCommentBehindLeftCurly is the optional one placed behind a left curly.
	InlineCommentBehindLeftCurly *parser.Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// InterpretExtend interprets *parser.Extend to *Extend.
func InterpretExtend(src *parser.Extend) (*Extend, error) {
	if src == nil {
		return nil, nil
	}

	extendBody, err := interpretExtendBody(src.ExtendBody)
	if err != nil {
		return nil, err
	}
	return &Extend{
		MessageType:                  src.MessageType,
		ExtendBody:                   extendBody,
		Comments:                     src.Comments,
		InlineComment:                src.InlineComment,
		InlineCommentBehindLeftCurly: src.InlineCommentBehindLeftCurly,
		Meta:                         src.Meta,
	}, nil
}

func interpretExtendBody(src []parser.Visitee) (
	*ExtendBody,
	error,
) {
	var fields []*parser.Field
	var emptyStatements []*parser.EmptyStatement
	for _, s := range src {
		switch t := s.(type) {
		case *parser.Field:
			fields = append(fields, t)
		case *parser.EmptyStatement:
			emptyStatements = append(emptyStatements, t)
		default:
			return nil, fmt.Errorf("invalid ExtendBody type %T of %v", t, t)
		}
	}
	return &ExtendBody{
		Fields:          fields,
		EmptyStatements: emptyStatements,
	}, nil
}
