package kmgTime

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

//approach 2
type ScratchPeriodList2 []ScratchPeriod2

func (p ScratchPeriodList2) Len() int {
	return len(p)
}
func (p ScratchPeriodList2) GetPeriodAtIndex(i int) Period {
	return p[i].Period
}
func (p ScratchPeriodList2) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type ScratchPeriod2 struct {
	Period
	ItemList []int //a list of item user can get at this period
}

func TestPeriodListInterface(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	periodList := ScratchPeriodList2{
		{
			Period:   Period{Start: MustFromMysqlFormat("2001-01-00 23:30:00"), End: MustFromMysqlFormat("2001-01-01 23:30:00")},
			ItemList: []int{1, 2},
		},
		{
			Period:   Period{Start: MustFromMysqlFormat("2001-01-03 23:30:00"), End: MustFromMysqlFormat("2001-01-04 23:30:00")},
			ItemList: []int{2, 3},
		},
		{
			Period:   Period{Start: MustFromMysqlFormat("2001-01-02 23:30:00"), End: MustFromMysqlFormat("2001-01-03 23:30:00")},
			ItemList: []int{3, 4},
		},
	}
	PeriodListSort(periodList)
	i, exist := SelectPeriodFromSortedPeriodList(MustFromMysqlFormat("2001-01-01 23:00:00"), periodList)
	t.Equal(exist, true)
	t.Equal(periodList[i].ItemList, []int{1, 2})
	i, exist = SelectPeriodFromSortedPeriodList(MustFromMysqlFormat("2001-01-03 23:00:00"), periodList)
	t.Equal(exist, true)
	t.Equal(periodList[i].ItemList, []int{3, 4})
}
