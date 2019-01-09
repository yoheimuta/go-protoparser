package meta

import "github.com/yoheimuta/go-protoparser/internal/lexer/scanner"

// Meta represents a meta information about the parsed element.
type Meta struct {
	// Pos is the source position.
	Pos Position
	// LastPos is the last source position.
	// Currently it is set when the parsed element type is message, enum, oneof, rpc or service.
	LastPos Position
}

// NewMeta creates a new Meta from scanner.Position.
func NewMeta(fromPos scanner.Position) Meta {
	return Meta{
		Pos: NewPosition(fromPos),
	}
}

// NewMetaWithLastPos creates a new Meta with LastPos from scanner.Position.
func NewMetaWithLastPos(
	fromPos scanner.Position,
	fromLastPos scanner.Position,
) Meta {
	return Meta{
		Pos:     NewPosition(fromPos),
		LastPos: NewPosition(fromLastPos),
	}
}
