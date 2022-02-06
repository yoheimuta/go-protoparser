package scanner

import (
	"unicode/utf8"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

// Position represents a source position.
type Position struct {
	meta.Position

	// columns is a map which the key is a line number and the value is a column number.
	columns map[int]int
}

// NewPosition creates a new Position.
func NewPosition() *Position {
	return &Position{
		Position: meta.Position{
			Offset: 0,
			Line:   1,
			Column: 1,
		},
		columns: make(map[int]int),
	}
}

// String stringify the position.
func (pos Position) String() string {
	return pos.Position.String()
}

// Advance advances the position value.
func (pos *Position) Advance(r rune) {
	length := utf8.RuneLen(r)
	pos.Offset += length

	if r == '\n' {
		pos.columns[pos.Line] = pos.Column

		pos.Line++
		pos.Column = 1
	} else {
		pos.Column++
	}
}

// AdvancedBulk returns a new position that advances the position value in a row.
func (pos Position) AdvancedBulk(s string) Position {
	for _, r := range s {
		pos.Advance(r)
	}
	last, _ := utf8.DecodeLastRuneInString(s)
	pos.Revert(last)
	return pos
}

// Revert reverts the position value.
func (pos *Position) Revert(r rune) {
	length := utf8.RuneLen(r)
	pos.Offset -= length

	if r == '\n' {
		pos.Line--
		pos.Column = pos.columns[pos.Line]
	} else {
		pos.Column--
	}
}
