package parser

import (
	"unicode/utf8"

	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// GroupField is one way to nest information in message definitions.
// proto2 only.
type GroupField struct {
	IsRepeated bool
	IsRequired bool
	IsOptional bool
	// GroupName must begin with capital letter.
	GroupName string
	// MessageBody can have fields, nested enum definitions, nested message definitions,
	// options, oneofs, map fields, extends, reserved, and extensions statements.
	MessageBody []Visitee
	FieldNumber string

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
func (f *GroupField) SetInlineComment(comment *Comment) {
	f.InlineComment = comment
}

// Accept dispatches the call to the visitor.
func (f *GroupField) Accept(v Visitor) {
	if !v.VisitGroupField(f) {
		return
	}

	for _, body := range f.MessageBody {
		body.Accept(v)
	}
	for _, comment := range f.Comments {
		comment.Accept(v)
	}
	if f.InlineComment != nil {
		f.InlineComment.Accept(v)
	}
	if f.InlineCommentBehindLeftCurly != nil {
		f.InlineCommentBehindLeftCurly.Accept(v)
	}
}

// ParseGroupField parses the group.
//  group = label "group" groupName "=" fieldNumber messageBody
//
// See https://developers.google.com/protocol-buffers/docs/reference/proto2-spec#group_field
func (p *Parser) ParseGroupField() (*GroupField, error) {
	var isRepeated bool
	var isRequired bool
	var isOptional bool
	p.lex.NextKeyword()
	startPos := p.lex.Pos

	if p.lex.Token == scanner.TREPEATED {
		isRepeated = true
	} else if p.lex.Token == scanner.TREQUIRED {
		isRequired = true
	} else if p.lex.Token == scanner.TOPTIONAL {
		isOptional = true
	} else {
		p.lex.UnNext()
	}

	p.lex.NextKeyword()
	if p.lex.Token != scanner.TGROUP {
		return nil, p.unexpected("group")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("groupName")
	}
	if !isCapitalized(p.lex.Text) {
		return nil, p.unexpectedf("groupName %q must begin with capital letter.", p.lex.Text)
	}
	groupName := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	fieldNumber, err := p.parseFieldNumber()
	if err != nil {
		return nil, p.unexpected("fieldNumber")
	}

	messageBody, inlineLeftCurly, lastPos, err := p.parseMessageBody()
	if err != nil {
		return nil, err
	}

	return &GroupField{
		IsRepeated:  isRepeated,
		IsRequired:  isRequired,
		IsOptional:  isOptional,
		GroupName:   groupName,
		FieldNumber: fieldNumber,
		MessageBody: messageBody,

		InlineCommentBehindLeftCurly: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

func (p *Parser) peekIsGroup() bool {
	p.lex.NextKeyword()
	switch p.lex.Token {
	case scanner.TREPEATED,
		scanner.TREQUIRED,
		scanner.TOPTIONAL:
		defer p.lex.UnNextTo(p.lex.RawText)
	default:
		p.lex.UnNext()
	}

	p.lex.NextKeyword()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TGROUP {
		return false
	}

	p.lex.Next()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TIDENT {
		return false
	}
	if !isCapitalized(p.lex.Text) {
		return false
	}

	p.lex.Next()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TEQUALS {
		return false
	}

	_, err := p.parseFieldNumber()
	defer p.lex.UnNextTo(p.lex.RawText)
	if err != nil {
		return false
	}

	p.lex.Next()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TLEFTCURLY {
		return false
	}
	return true
}

// isCapitalized returns true if is not empty and the first letter is
// an uppercase character.
func isCapitalized(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	return isUpper(r)
}

func isUpper(r rune) bool {
	return 'A' <= r && r <= 'Z'
}
