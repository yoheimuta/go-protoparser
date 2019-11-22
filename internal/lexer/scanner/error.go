package scanner

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/issues"
)

func (s *Scanner) unexpected(found rune, expected string) issues.ParseError {
	_, file, line, _ := runtime.Caller(1)
	pe := issues.ParseError{
		Filename: s.pos.Filename,
		Line:     s.pos.Line,
		Column:   s.pos.Column,
		Found:    string(found),
		Expected: expected,
	}
	pe.SetOccured(file, line)
	return pe
}
