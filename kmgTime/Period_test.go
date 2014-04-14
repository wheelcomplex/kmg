package kmgTime

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
	"time"
)

func TestGetPeriodFromSortedList(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	SortedList := []Period{
		{
			Start: MustFromMysqlFormat("2001-01-01 00:00:00"),
			End:   MustFromMysqlFormat("2001-01-01 01:00:00"),
		},
		{
			Start: MustFromMysqlFormat("2001-01-01 02:00:00"),
			End:   MustFromMysqlFormat("2001-01-01 03:00:00"),
		},
		{
			Start: MustFromMysqlFormat("2001-01-01 03:00:00"),
			End:   MustFromMysqlFormat("2001-01-01 04:00:00"),
		},
	}
	for _, testcase := range []struct {
		t  time.Time
		i  int
		ok bool
	}{
		{MustFromMysqlFormat("2001-01-00 23:30:00"), 0, false},
		{MustFromMysqlFormat("2001-01-01 00:30:00"), 0, true},
		{MustFromMysqlFormat("2001-01-01 03:00:00"), 2, true},
		{MustFromMysqlFormat("2001-01-01 04:30:00"), 0, false},
	} {
		i, ok := GetPeriodFromSortedList(testcase.t, SortedList)
		t.Equal(i, testcase.i)
		t.Equal(ok, testcase.ok)
	}
}

type ScratchPeriodList []ScratchPeriod

func (p ScratchPeriodList) Len() int {
	return len(p)
}
func (p ScratchPeriodList) Less(i, j int) bool {
	return p[i].Period.Start.Before(p[j].Period.Start)
}
func (p ScratchPeriodList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p ScratchPeriodList) GetPeriodAtIndex(i int) Period {
	return p[i].Period
}

type ScratchPeriod struct {
	Period
	ItemList []int
}

func (p ScratchPeriod) GetPeriod() Period {
	return p.Period
}

//approach 1
func TestGetPeriodFromGenericSortedList(ot *testing.T) {
	/*
		t := kmgTest.NewTestTools(ot)
		PeriodList:=
		SortedList := []Period{
			{
				Start: MustFromMysqlFormat("2001-01-01 00:00:00"),
				End:   MustFromMysqlFormat("2001-01-01 01:00:00"),
			},
			{
				Start: MustFromMysqlFormat("2001-01-01 02:00:00"),
				End:   MustFromMysqlFormat("2001-01-01 03:00:00"),
			},
			{
				Start: MustFromMysqlFormat("2001-01-01 03:00:00"),
				End:   MustFromMysqlFormat("2001-01-01 04:00:00"),
			},
		}
		for _, testcase := range []struct {
			t  time.Time
			i  int
			ok bool
		}{
			{MustFromMysqlFormat("2001-01-00 23:30:00"), 0, false},
			{MustFromMysqlFormat("2001-01-01 00:30:00"), 0, true},
			{MustFromMysqlFormat("2001-01-01 03:00:00"), 2, true},
			{MustFromMysqlFormat("2001-01-01 04:30:00"), 0, false},
		} {
			i, ok := GetPeriodFromSortedList(testcase.t, SortedList)
			t.Equal(i, testcase.i)
			t.Equal(ok, testcase.ok)
		}
	*/
}
