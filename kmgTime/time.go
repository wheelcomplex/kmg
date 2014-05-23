package kmgTime

import (
	"fmt"
	"time"
)

const (
	FormatMysql     = "2006-01-02 15:04:05"
	FormatDateMysql = "2006-01-02"
	Iso3339Hour     = "2006-01-02T15"
	Iso3339Minute   = "2006-01-02T15:04"
	Iso3339Second   = "2006-01-02T15:04:05"
	AppleJsonFormat = "2006-01-02 15:04:05 Etc/MST" //仅解决GMT的这个特殊情况.其他不管,如果苹果返回的字符串换时区了就悲剧了
	Day             = 24 * time.Hour
)

var ParseFormatGuessList = []string{
	FormatMysql,
	FormatDateMysql,
	Iso3339Hour,
	Iso3339Minute,
	Iso3339Second,
}

var BeijingZone = time.FixedZone("CST", 8*60*60)
var ESTZone = time.FixedZone("EST", -5*60*60) // Eastern Standard Time(加拿大)

func ParseAutoInLocal(sTime string) (t time.Time, err error) {
	return ParseAutoInLocation(sTime, time.Local)
}

func MustParseAutoInLocal(sTime string) (t time.Time) {
	t, err := ParseAutoInLocation(sTime, time.Local)
	if err != nil {
		panic(err)
	}
	return t
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
	return t.Local()
}

func ToDateString(t time.Time) string {
	return t.Format(FormatDateMysql)
}

//规则到日期,去掉时分秒
func ToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

//规整到日期然后相减
func DateSub(t1 time.Time, t2 time.Time, loc *time.Location) time.Duration {
	return ToDate(t1.In(loc)).Sub(ToDate(t2.In(loc)))
}

func DateSubLocal(t1 time.Time, t2 time.Time) time.Duration {
	return DateSub(t1, t2, time.Local)
}
