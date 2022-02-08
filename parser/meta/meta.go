package meta

// Meta represents a meta information about the parsed element.
type Meta struct {
	// Pos is the source position.
	Pos Position
	// LastPos is the last source position.
	// Currently it is set when the parsed element type is
	// syntax, package, comment, import, option, message, enum, oneof, rpc or service.
	LastPos Position
}
