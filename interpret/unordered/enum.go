package unordered

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

// EnumBody is unordered in nature, but each slice field preserves the original order.
type EnumBody struct {
	Options         []*parser.Option
	EnumFields      []*parser.EnumField
	EmptyStatements []*parser.EmptyStatement
}

// Enum consists of a name and an enum body.
type Enum struct {
	EnumName string
	EnumBody *EnumBody

	// Comments are the optional ones placed at the beginning.
	Comments []*parser.Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *parser.Comment
	// InlineCommentBehindLeftCurly is the optional one placed behind a left curly.
	InlineCommentBehindLeftCurly *parser.Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// InterpretEnum interprets *parser.Enum to *Enum.
func InterpretEnum(src *parser.Enum) (*Enum, error) {
	if src == nil {
		return nil, nil
	}

	enumBody, err := interpretEnumBody(src.EnumBody)
	if err != nil {
		return nil, err
	}
	return &Enum{
		EnumName:                     src.EnumName,
		EnumBody:                     enumBody,
		Comments:                     src.Comments,
		InlineComment:                src.InlineComment,
		InlineCommentBehindLeftCurly: src.InlineCommentBehindLeftCurly,
		Meta:                         src.Meta,
	}, nil
}

func interpretEnumBody(src []parser.Visitee) (
	*EnumBody,
	error,
) {
	var options []*parser.Option
	var enumFields []*parser.EnumField
	var emptyStatements []*parser.EmptyStatement
	for _, s := range src {
		switch t := s.(type) {
		case *parser.Option:
			options = append(options, t)
		case *parser.EnumField:
			enumFields = append(enumFields, t)
		case *parser.EmptyStatement:
			emptyStatements = append(emptyStatements, t)
		default:
			return nil, fmt.Errorf("invalid EnumBody type %v of %v", t, s)
		}
	}
	return &EnumBody{
		Options:         options,
		EnumFields:      enumFields,
		EmptyStatements: emptyStatements,
	}, nil
}
