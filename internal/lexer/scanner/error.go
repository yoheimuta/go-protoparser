package scanner

import (
	"fmt"
	"runtime"
)

func (s *Scanner) unexpected(found rune, expected string) error {
	_, file, line, _ := runtime.Caller(1)
	message := fmt.Sprintf(" at %s:%d", file, line)
	return fmt.Errorf("found %q but expected [%s]%s", found, expected, message)
}
