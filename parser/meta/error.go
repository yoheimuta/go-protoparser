// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package meta

import "fmt"

// Error is the error type returned for all scanning/lexing/parsing related errors.
type Error struct {
	Pos      Position
	Expected string
	Found    string

	occuredIn string
	occuredAt int
}

func (e *Error) Error() string {
	if e.occuredAt == 0 && e.occuredIn == "" {
		return fmt.Sprintf("found %q but expected [%s]", e.Found, e.Expected)
	}
	return fmt.Sprintf("found %q but expected [%s] at %s:%d", e.Found, e.Expected, e.occuredIn, e.occuredAt)
}

// SetOccured sets the file and the line number at which the error was raised (through runtime.Caller).
func (e *Error) SetOccured(occuredIn string, occuredAt int) {
	e.occuredIn = occuredIn
	e.occuredAt = occuredAt
}
