package scanner

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/errors"
)

func (s *Scanner) unexpected(found rune, expected string) errors.ParseError {
	_, file, line, _ := runtime.Caller(1)
	return errors.NewParseError(string(found), expected, file, line)
}
