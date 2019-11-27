// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package meta

import "fmt"

type Error struct {
	Pos      Position
	Expected string
	Found    string

	occuredIn string
	occuredAt int
}

func (e *Error) Error() string {
	if e.occuredAt == 0 && e.occuredIn == "" {
		return fmt.Sprintf("found %s but expected [%s]", e.Found, e.Expected)
	}
	return fmt.Sprintf("found %s but expected [%s] at %s:%d", e.Found, e.Expected, e.occuredIn, e.occuredAt)
}

func (e *Error) SetOccured(occuredIn string, occuredAt int) {
	e.occuredIn = occuredIn
	e.occuredAt = occuredAt
}
