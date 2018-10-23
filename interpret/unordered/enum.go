package unordered

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/parser"
)

// Enum consists of a name and an enum body.
type Enum struct {
	EnumName string

	// EnumBody is unordered in nature, but each slice field preserves the original order.
	Options         []*parser.Option
	EnumFields      []*parser.EnumField
	EmptyStatements []*parser.EmptyStatement

	// Comments are the optional ones placed at the beginning.
	Comments []*parser.Comment
}

// InterpretEnum interprets *parser.Enum to *Enum.
func InterpretEnum(src *parser.Enum) (*Enum, error) {
	if src == nil {
		return nil, nil
	}

	options, enumFields, emptyStatements, err := interpretEnumBody(src.EnumBody)
	if err != nil {
		return nil, err
	}
	return &Enum{
		EnumName:        src.EnumName,
		Options:         options,
		EnumFields:      enumFields,
		EmptyStatements: emptyStatements,
		Comments:        src.Comments,
	}, nil
}

func interpretEnumBody(src []interface{}) (
	options []*parser.Option,
	enumFields []*parser.EnumField,
	emptyStatements []*parser.EmptyStatement,
	err error,
) {
	for _, s := range src {
		switch t := s.(type) {
		case *parser.Option:
			options = append(options, t)
		case *parser.EnumField:
			enumFields = append(enumFields, t)
		case *parser.EmptyStatement:
			emptyStatements = append(emptyStatements, t)
		default:
			err = fmt.Errorf("invalid EnumBody type %v of %v", t, s)
		}
	}
	return
}
