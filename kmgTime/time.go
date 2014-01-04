package kmgTime

import (
	"fmt"
	"time"
)

const (
	FormatMysql   = "2006-01-02 15:04:05"
	Iso3339Hour   = "2006-01-02T15"
	Iso3339Minute = "2006-01-02T15:04"
	Iso3339Second = "2006-01-02T15:04:05"
)

var ParseFormatGuessList = []string{
	FormatMysql,
	Iso3339Hour,
	Iso3339Minute,
	Iso3339Second,
}

func ParseAutoInLocal(sTime string) (t time.Time, err error) {
	return ParseAutoInLocation(sTime, time.Local)
}

//auto guess format from ParseFormatGuessList
func ParseAutoInLocation(sTime string, loc *time.Location) (t time.Time, err error) {
	for _, format := range ParseFormatGuessList {
		t, err = time.ParseInLocation(format, sTime, loc)
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("[ParseAutoInLocation] time: %s can not parse", sTime)
	return
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
