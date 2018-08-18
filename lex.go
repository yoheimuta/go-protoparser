package protoparser

import (
	"io"
	"log"
	"path/filepath"
	"runtime"
	"text/scanner"
)

type lexer struct {
	scan  scanner.Scanner
	token rune

	debug bool
}

type option func(*lexer)

func withDebug(debug bool) option {
	return func(l *lexer) {
		l.debug = debug
	}
}

func newlexer(input io.Reader, opts ...option) *lexer {
	lex := new(lexer)
	for _, opt := range opts {
		opt(lex)
	}

	lex.scan.Init(input)
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanComments
	lex.next()
	return lex
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()

	if lex.debug {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf(
				"[DEBUG] Token : [%s], position [%v] called from %s:%d\n",
				lex.text(),
				lex.scan.Pos(),
				filepath.Base(file),
				line,
			)
		}
	}
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}
