package lexer

import (
	"io"
	"log"

	"path/filepath"
	"runtime"

	"github.com/yoheimuta/go-protoparser/internal/lexer/scanner"
)

// Lexer2 is a lexer.
type Lexer2 struct {
	// Token is the lexical token.
	Token scanner.Token

	// Text is the lexical value.
	Text string

	// Error is called for each error encountered. If no Error
	// function is set, the error is reported to os.Stderr.
	Error func(lexer *Lexer2, err error)

	scanner *scanner.Scanner
	scanErr error
	debug   bool
}

// Option2 is an option for lexer.NewLexer2.
type Option2 func(*Lexer2)

// WithDebug2 is an option to enable the debug mode.
func WithDebug2(debug bool) Option2 {
	return func(l *Lexer2) {
		l.debug = debug
	}
}

// NewLexer2 creates a new lexer.
func NewLexer2(input io.Reader, opts ...Option2) *Lexer2 {
	lex := new(Lexer2)
	for _, opt := range opts {
		opt(lex)
	}

	lex.Error = func(_ *Lexer2, err error) {
		log.Printf(`Lexer encountered the error "%v"`, err)
	}
	lex.scanner = scanner.NewScanner(input)
	return lex
}

// Next scans the read buffer.
func (lex *Lexer2) Next() {
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
func (lex *Lexer2) NextKeywordOrStrLit() {
	lex.nextWithSpecificMode(scanner.ScanKeyword | scanner.ScanStrLit)
}

// NextKeyword scans the read buffer with ScanKeyword mode.
func (lex *Lexer2) NextKeyword() {
	lex.nextWithSpecificMode(scanner.ScanKeyword)
}

// NextStrLit scans the read buffer with ScanStrLit mode.
func (lex *Lexer2) NextStrLit() {
	lex.nextWithSpecificMode(scanner.ScanStrLit)
}

func (lex *Lexer2) nextWithSpecificMode(nextMode scanner.Mode) {
	mode := lex.scanner.Mode
	defer func() {
		lex.scanner.Mode = mode
	}()

	lex.scanner.Mode = nextMode
	lex.Next()
}

// IsEOF checks whether read buffer is empty.
func (lex *Lexer2) IsEOF() bool {
	return lex.Token == scanner.TEOF
}

// LatestErr returns the latest non-EOF error that was encountered by the Lexer2.Next().
func (lex *Lexer2) LatestErr() error {
	return lex.scanErr
}

// UnNext put the latest text back to the read buffer.
func (lex *Lexer2) UnNext() {
	lex.scanner.UnScan(lex.Text)
	lex.Token = scanner.TILLEGAL
}
