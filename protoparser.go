package protoparser

import (
	"io"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/parser"
)

// ParseConfig is a config for parser.
type ParseConfig struct {
	debug      bool
	permissive bool
}

// Option is an option for ParseConfig.
type Option func(*ParseConfig)

// WithDebug is an option to enable the debug mode.
func WithDebug(debug bool) Option {
	return func(l *ParseConfig) {
		l.debug = debug
	}
}

// WithPermissive is an option to allow the permissive parsing rather than the just documented spec.
func WithPermissive(permissive bool) Option {
	return func(l *ParseConfig) {
		l.permissive = permissive
	}
}

// Parse parses a Protocol Buffer file.
func Parse(input io.Reader, options ...Option) (*parser.Proto, error) {
	config := &ParseConfig{}
	for _, opt := range options {
		opt(config)
	}

	p := parser.NewParser(lexer.NewLexer2(input, lexer.WithDebug2(config.debug), lexer.WithPermissive(config.permissive)))
	return p.ParseProto()
}
