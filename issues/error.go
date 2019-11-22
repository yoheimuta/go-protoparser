// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package issues

import (
	"fmt"
)

// ParseError is the error returned during parsing.
type ParseError struct {
	Filename string
	Line     int
	Column   int

	Found      string
	Expected   string
	occurredIn string
	occurredAt int
}

func (pe ParseError) String() string {
	out := fmt.Sprintf("found %q but expected [%s]", pe.Found, pe.Expected)

	if pe.occurredIn != "" && pe.occurredAt != 0 {
		out += fmt.Sprintf(" at %s:%d", pe.occurredIn, pe.occurredAt)
	}

	return out
}

func (pe ParseError) Error() string {
	return pe.String()
}

// SetOccurred sets the fields to log where each error was raised.
func (pe *ParseError) SetOccured(occurerdIn string, occurredAt int) {
	pe.occurredIn = occurerdIn
	pe.occurredAt = occurredAt
}
