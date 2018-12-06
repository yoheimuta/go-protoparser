package scanner

import (
	"fmt"
	"unicode/utf8"
)

// Position represents a source position.
type Position struct {
	// Offset is a byte offset, starting at 0
	Offset int
	// Line is a line number, starting at 1
	Line int
	// Column is a column number, starting at 1 (character count per line)
	Column int

	// columns is a map which the key is a line number and the value is a column number.
	columns map[int]int
}

// NewPosition creates a new Position.
func NewPosition() *Position {
	return &Position{
		Offset:  0,
		Line:    1,
		Column:  1,
		columns: make(map[int]int),
	}
}

// String stringify the position.
func (pos Position) String() string {
	s := "<input>"
	s += fmt.Sprintf(":%d:%d", pos.Line, pos.Column)
	return s
}

// Advance advances the position value.
func (pos *Position) Advance(r rune) {
	len := utf8.RuneLen(r)
	pos.Offset += len

	if r == '\n' {
		pos.columns[pos.Line] = pos.Column

		pos.Line++
		pos.Column = 1
	} else {
		pos.Column++
	}
}

// Revert reverts the position value.
func (pos *Position) Revert(r rune) {
	len := utf8.RuneLen(r)
	pos.Offset -= len

	if r == '\n' {
		pos.Line--
		pos.Column = pos.columns[pos.Line]
	} else {
		pos.Column--
	}
}
