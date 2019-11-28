package scanner

import (
	"runtime"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func (s *Scanner) unexpected(found rune, expected string) error {
	_, file, line, _ := runtime.Caller(1)
	err := &meta.Error{
		Pos:      s.pos.Position,
		Expected: expected,
		Found:    string(found),
	}
	err.SetOccured(file, line)
	return err
}
