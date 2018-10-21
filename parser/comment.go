package parser

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

const (
	cStyleCommentPrefix     = "/*"
	cStyleCommentSuffix     = "*/"
	cPlusStyleCommentPrefix = "//"
)

// Comment is a comment in either C/C++-style // and /* ... */ syntax.
type Comment struct {
	// Raw includes a comment syntax like // and /* */.
	Raw string
}

// IsCStyle refers to /* ... */.
func (c *Comment) IsCStyle() bool {
	return strings.HasPrefix(c.Raw, cStyleCommentPrefix)
}

// Lines formats comment text lines without prefixes //, /* or suffix */.
func (c *Comment) Lines() []string {
	raw := c.Raw
	if c.IsCStyle() {
		raw = strings.TrimPrefix(raw, cStyleCommentPrefix)
		raw = strings.TrimSuffix(raw, cStyleCommentSuffix)
	} else {
		raw = strings.TrimPrefix(raw, cPlusStyleCommentPrefix)
	}
	return strings.Split(raw, "\n")
}

// ParseComments parsers a sequence of comments.
//  comments = { comment }
//
// See https://developers.google.com/protocol-buffers/docs/proto3#adding-comments
func (p *Parser) ParseComments() []*Comment {
	var comments []*Comment
	for {
		comment, err := p.parseComment()
		if err != nil {
			// ignores the err because the comment is optional.
			return comments
		}
		comments = append(comments, comment)
	}
}

// See https://developers.google.com/protocol-buffers/docs/proto3#adding-comments
func (p *Parser) parseComment() (*Comment, error) {
	p.lex.NextComment()
	if p.lex.Token == scanner.TCOMMENT {
		return &Comment{
			Raw: p.lex.Text,
		}, nil
	}
	defer p.lex.UnNext()
	return nil, p.unexpected("comment")
}
