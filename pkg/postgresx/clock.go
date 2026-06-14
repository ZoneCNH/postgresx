package postgresx

import "time"

// Clock abstracts time for tests and deterministic health output.
type Clock interface {
	Now() time.Time
}

// RealClock returns the current wall-clock time.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// NewRealClock creates a production clock.
func NewRealClock() Clock { return RealClock{} }
