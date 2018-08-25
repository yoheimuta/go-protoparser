package lexer

import (
	"io"
	"log"
	"path/filepath"
	"runtime"
	"text/scanner"
)

// Lexer is a lexer.
type Lexer struct {
	// Token is the text/scanner token.
	Token rune

	scan  scanner.Scanner
	debug bool
}

// Option is an option for lexer.NewLexer.
type Option func(*Lexer)

// WithDebug is an option to enable the debug mode.
func WithDebug(debug bool) Option {
	return func(l *Lexer) {
		l.debug = debug
	}
}

// NewLexer creates a new lexer.
func NewLexer(input io.Reader, opts ...Option) *Lexer {
	lex := new(Lexer)
	for _, opt := range opts {
		opt(lex)
	}

	lex.scan.Init(input)
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanComments
	lex.Next()
	return lex
}

// Next scans the internal buffer.
func (lex *Lexer) Next() {
	lex.Token = lex.scan.Scan()

	if lex.debug {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf(
				"[DEBUG] Token : [%s], position [%v] called from %s:%d\n",
				lex.Text(),
				lex.scan.Pos(),
				filepath.Base(file),
				line,
			)
		}
	}
}

// Text returns the current Token text.
func (lex *Lexer) Text() string {
	return lex.scan.TokenText()
}
