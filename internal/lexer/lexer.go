package lexer

import (
	"io"
	"log"

	"path/filepath"
	"runtime"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

// Lexer is a lexer.
type Lexer struct {
	// Token is the lexical token.
	Token scanner.Token

	// Text is the lexical value.
	Text string

	// Error is called for each error encountered. If no Error
	// function is set, the error is reported to os.Stderr.
	Error func(lexer *Lexer, err error)

	scanner *scanner.Scanner
	scanErr error
	debug   bool
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

	lex.Error = func(_ *Lexer, err error) {
		log.Printf(`Lexer encountered the error "%v"`, err)
	}
	lex.scanner = scanner.NewScanner(input)
	return lex
}

// Next scans the read buffer.
func (lex *Lexer) Next() {
	defer func() {
		if lex.debug {
			_, file, line, ok := runtime.Caller(2)
			if ok {
				log.Printf(
					"[DEBUG] Text=[%s], Token=[%v] called from %s:%d\n",
					lex.Text,
					lex.Token,
					filepath.Base(file),
					line,
				)
			}
		}
	}()

	var err error
	lex.Token, lex.Text, err = lex.scanner.Scan()
	if err != nil {
		lex.scanErr = err
		lex.Error(lex, err)
	}
}

// NextKeywordOrStrLit scans the read buffer with ScanKeyword or ScanStrLit modes.
func (lex *Lexer) NextKeywordOrStrLit() {
	lex.nextWithSpecificMode(scanner.ScanKeyword | scanner.ScanStrLit)
}

// NextKeyword scans the read buffer with ScanKeyword mode.
func (lex *Lexer) NextKeyword() {
	lex.nextWithSpecificMode(scanner.ScanKeyword)
}

// NextStrLit scans the read buffer with ScanStrLit mode.
func (lex *Lexer) NextStrLit() {
	lex.nextWithSpecificMode(scanner.ScanStrLit)
}

// NextLit scans the read buffer with ScanLit mode.
func (lex *Lexer) NextLit() {
	lex.nextWithSpecificMode(scanner.ScanLit)
}

// NextNumberLit scans the read buffer with ScanNumberLit mode.
func (lex *Lexer) NextNumberLit() {
	lex.nextWithSpecificMode(scanner.ScanNumberLit)
}

func (lex *Lexer) nextWithSpecificMode(nextMode scanner.Mode) {
	mode := lex.scanner.Mode
	defer func() {
		lex.scanner.Mode = mode
	}()

	lex.scanner.Mode = nextMode
	lex.Next()
}

// IsEOF checks whether read buffer is empty.
func (lex *Lexer) IsEOF() bool {
	return lex.Token == scanner.TEOF
}

// LatestErr returns the latest non-EOF error that was encountered by the Lexer.Next().
func (lex *Lexer) LatestErr() error {
	return lex.scanErr
}

// UnNext put the latest text back to the read buffer.
func (lex *Lexer) UnNext() {
	lex.scanner.UnScan(lex.Text)
	lex.Token = scanner.TILLEGAL
}
