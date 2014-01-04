package kmgTime

import "time"

type Nower interface {
	Now() time.Time
}

func GetDefaultNower() Nower {
	return DefaultNower
}

var DefaultNower tDefaultNower

func NewFixedNower(time time.Time) Nower {
	return FixedNower{time}
}

type tDefaultNower struct{}

func (nower tDefaultNower) Now() time.Time {
	return time.Now()
}

type FixedNower struct {
	Time time.Time
}

func (nower FixedNower) Now() time.Time {
	return nower.Time
}
