package kmgTime

import "time"

const (
	FormatMysql = "2006-01-02 15:04:05"
)

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

//utc time
func MustFromMysqlFormat(timeString string) time.Time {
	t, err := time.Parse(FormatMysql, timeString)
	if err != nil {
		panic(err)
	}
	return t
}

//local time
func MustFromLocalMysqlFormat(timeString string) time.Time {
	t, err := time.ParseInLocation(FormatMysql, timeString, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func ToLocal(t time.Time) time.Time {
	return t.In(time.Local)
}
