package parser

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

type parseReservedErr struct {
	parseRangesErr     error
	parseFieldNamesErr error
}

func (e *parseReservedErr) Error() string {
	return fmt.Sprintf("%v:%v", e.parseRangesErr, e.parseFieldNamesErr)
}

// Range is a range of field numbers. End is an optional value.
type Range struct {
	Begin string
	End   string
}

// Reserved declares a range of field numbers or field names that cannot be used in this message.
// These component Ranges and FieldNames are mutually exclusive.
type Reserved struct {
	Ranges     []*Range
	FieldNames []string

	// Comments are the optional ones placed at the beginning.
	Comments []*Comment
	// InlineComment is the optional one placed at the ending.
	InlineComment *Comment
	// Meta is the meta information.
	Meta meta.Meta
}

// SetInlineComment implements the HasInlineCommentSetter interface.
func (r *Reserved) SetInlineComment(comment *Comment) {
	r.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (r *Reserved) Accept(v Visitor) {
	if !v.VisitReserved(r) {
		return
	}

	for _, comment := range r.Comments {
		comment.Accept(v)
	}
	if r.InlineComment != nil {
		r.InlineComment.Accept(v)
	}
}

// ParseReserved parses the reserved.
//  reserved = "reserved" ( ranges | fieldNames ) ";"
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#reserved
func (p *Parser) ParseReserved() (*Reserved, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TRESERVED {
		return nil, p.unexpected("reserved")
	}
	startPos := p.lex.Pos

	parse := func() ([]*Range, []string, error) {
		ranges, err := p.parseRanges()
		if err == nil {
			return ranges, nil, nil
		}

		fieldNames, ferr := p.parseFieldNames()
		if ferr == nil {
			return nil, fieldNames, nil
		}

		return nil, nil, &parseReservedErr{
			parseRangesErr:     err,
			parseFieldNamesErr: ferr,
		}
	}

	ranges, fieldNames, err := parse()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &Reserved{
		Ranges:     ranges,
		FieldNames: fieldNames,
		Meta:       meta.Meta{Pos: startPos.Position},
	}, nil
}

// ranges = range { "," range }
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#reserved
func (p *Parser) parseRanges() ([]*Range, error) {
	var ranges []*Range
	rangeValue, err := p.parseRange()
	if err != nil {
		return nil, err
	}
	ranges = append(ranges, rangeValue)

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}

		rangeValue, err := p.parseRange()
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, rangeValue)
	}
	return ranges, nil
}

// range =  intLit [ "to" ( intLit | "max" ) ]
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#reserved
func (p *Parser) parseRange() (*Range, error) {
	p.lex.NextNumberLit()
	if p.lex.Token != scanner.TINTLIT {
		p.lex.UnNext()
		return nil, p.unexpected("intLit")
	}
	begin := p.lex.Text

	p.lex.Next()
	if p.lex.Text != "to" {
		p.lex.UnNext()
		return &Range{
			Begin: begin,
		}, nil
	}

	p.lex.NextNumberLit()
	switch {
	case p.lex.Token == scanner.TINTLIT,
		p.lex.Text == "max":
		return &Range{
			Begin: begin,
			End:   p.lex.Text,
		}, nil
	default:
		break
	}
	return nil, p.unexpected(`"intLit | "max"`)
}

// fieldNames = fieldName { "," fieldName }
// See https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#reserved
func (p *Parser) parseFieldNames() ([]string, error) {
	var fieldNames []string

	fieldName, err := p.parseQuotedFieldName()
	if err != nil {
		return nil, err
	}
	fieldNames = append(fieldNames, fieldName)

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}

		fieldName, err = p.parseQuotedFieldName()
		if err != nil {
			return nil, err
		}
		fieldNames = append(fieldNames, fieldName)
	}
	return fieldNames, nil
}

// quotedFieldName = quote + fieldName + quote
// TODO: Fixed according to defined documentation. Currently(2018.10.16) the reference lacks the spec.
// See https://github.com/protocolbuffers/protobuf/issues/4558
func (p *Parser) parseQuotedFieldName() (string, error) {
	p.lex.NextStrLit()
	if p.lex.Token != scanner.TSTRLIT {
		p.lex.UnNext()
		return "", p.unexpected("quotedFieldName")
	}
	return p.lex.Text, nil
}
