package protoparser

import (
	"io"

	"github.com/yoheimuta/go-protoparser/v4/interpret/unordered"
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

// ParseConfig is a config for parser.
type ParseConfig struct {
	debug                 bool
	permissive            bool
	bodyIncludingComments bool
	filename              string
}

// Option is an option for ParseConfig.
type Option func(*ParseConfig)

// WithDebug is an option to enable the debug mode.
func WithDebug(debug bool) Option {
	return func(c *ParseConfig) {
		c.debug = debug
	}
}

// WithPermissive is an option to allow the permissive parsing rather than the just documented spec.
func WithPermissive(permissive bool) Option {
	return func(c *ParseConfig) {
		c.permissive = permissive
	}
}

// WithBodyIncludingComments is an option to allow to include comments into each element's body.
// The comments are remaining of other elements'Comments and InlineComment.
func WithBodyIncludingComments(bodyIncludingComments bool) Option {
	return func(c *ParseConfig) {
		c.bodyIncludingComments = bodyIncludingComments
	}
}

// WithFilename is an option to set filename to the Position.
func WithFilename(filename string) Option {
	return func(c *ParseConfig) {
		c.filename = filename
	}
}

// Parse parses a Protocol Buffer file.
func Parse(input io.Reader, options ...Option) (*parser.Proto, error) {
	config := &ParseConfig{
		permissive: true,
	}
	for _, opt := range options {
		opt(config)
	}

	p := parser.NewParser(
		lexer.NewLexer(
			input,
			lexer.WithDebug(config.debug),
			lexer.WithFilename(config.filename),
		),
		parser.WithPermissive(config.permissive),
		parser.WithBodyIncludingComments(config.bodyIncludingComments),
	)
	return p.ParseProto()
}

// UnorderedInterpret interprets a Proto to an unordered one without interface{}.
func UnorderedInterpret(proto *parser.Proto) (*unordered.Proto, error) {
	return unordered.InterpretProto(proto)
}
