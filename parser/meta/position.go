package meta

import (
	"fmt"
)

// Position represents a source position.
type Position struct {
	// Filename is a name of file, if any
	Filename string
	// Offset is a byte offset, starting at 0
	Offset int
	// Line is a line number, starting at 1
	Line int
	// Column is a column number, starting at 1 (character count per line)
	Column int
}

// String stringify the position.
func (pos Position) String() string {
	s := pos.Filename
	if s == "" {
		s = "<input>"
	}
	s += fmt.Sprintf(":%d:%d", pos.Line, pos.Column)
	return s
}
