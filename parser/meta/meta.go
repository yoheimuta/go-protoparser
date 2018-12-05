package meta

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// Meta represents a meta information about the parsed element.
type Meta struct {
	// Pos is the source position.
	Pos Position
}

// NewMeta creates a new Meta from scanner.Position.
func NewMeta(fromPos scanner.Position) Meta {
	return Meta{
		Pos: NewPosition(fromPos),
	}
}
