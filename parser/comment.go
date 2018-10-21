package parser

import "strings"

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
