package lexer

import (
	"fmt"
	"runtime"
)

func (lex *Lexer2) unexpected(found, expected string) error {
	debug := ""
	if lex.debug {
		_, file, line, _ := runtime.Caller(1)
		debug = fmt.Sprintf(" at %s:%d", file, line)
	}
	return fmt.Errorf("found %q but expected [%s]%s", found, expected, debug)
}
