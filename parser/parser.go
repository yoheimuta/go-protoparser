package parser

import "github.com/yoheimuta/go-protoparser/v4/lexer"

// Parser is a parser.
type Parser struct {
	lex *lexer.Lexer

	permissive            bool
	bodyIncludingComments bool
}

// ConfigOption is an option for Parser.
type ConfigOption func(*Parser)

// WithPermissive is an option to allow the permissive parsing rather than the just documented spec.
func WithPermissive(permissive bool) ConfigOption {
	return func(p *Parser) {
		p.permissive = permissive
	}
}

// WithBodyIncludingComments is an option to allow to include comments into each element's body.
// The comments are remaining of other elements'Comments and InlineComment.
func WithBodyIncludingComments(bodyIncludingComments bool) ConfigOption {
	return func(p *Parser) {
		p.bodyIncludingComments = bodyIncludingComments
	}
}

// NewParser creates a new Parser.
func NewParser(lex *lexer.Lexer, opts ...ConfigOption) *Parser {
	p := &Parser{
		lex: lex,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// IsEOF checks whether the lex's read buffer is empty.
func (p *Parser) IsEOF() bool {
	p.lex.Next()
	defer p.lex.UnNext()
	return p.lex.IsEOF()
}
