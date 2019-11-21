package scanner

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/issues"
)

func (s *Scanner) unexpected(found rune, expected string) issues.ParseError {
	_, file, line, _ := runtime.Caller(1)
	return issues.NewParseError(string(found), expected, file, line)
}
