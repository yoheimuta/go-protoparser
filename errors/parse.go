// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package errors

import (
	"fmt"
)

// ParseError is the error returned during parsing.
type ParseError struct {
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

func NewParseError(found string, expected string, occuredIn string, occuredAt int) ParseError {
	return ParseError{
		Found:      found,
		Expected:   expected,
		occurredIn: occuredIn,
		occurredAt: occuredAt,
	}
}
