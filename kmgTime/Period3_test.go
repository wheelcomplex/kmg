package kmgTime

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

//approach 3
type ScratchItemList3 [][]int

func TestPeriodList(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	itemList := ScratchItemList3{
		[]int{1, 2},
		[]int{2, 3},
		[]int{3, 4},
	}
	periodList := PeriodList{
		{
			Period:      Period{Start: MustFromMysqlFormat("2001-01-00 23:30:00"), End: MustFromMysqlFormat("2001-01-01 23:30:00")},
			OriginIndex: 0,
		},
		{
			Period:      Period{Start: MustFromMysqlFormat("2001-01-03 23:30:00"), End: MustFromMysqlFormat("2001-01-04 23:30:00")},
			OriginIndex: 1,
		},
		{
			Period:      Period{Start: MustFromMysqlFormat("2001-01-02 23:30:00"), End: MustFromMysqlFormat("2001-01-03 23:30:00")},
			OriginIndex: 2,
		},
	}
	periodList.Sort()
	i, exist := periodList.SelectFromTime(MustFromMysqlFormat("2001-01-01 23:00:00"))
	t.Equal(exist, true)
	t.Equal(itemList[i], []int{1, 2})
	i, exist = periodList.SelectFromTime(MustFromMysqlFormat("2001-01-03 23:00:00"))
	t.Equal(exist, true)
	t.Equal(itemList[i], []int{3, 4})
}
