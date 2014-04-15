package kmgTime

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
	"sort"
)

//approach 1
type ScratchPeriodList1 []ScratchPeriod1

func (p ScratchPeriodList1) Len() int {
	return len(p)
}
func (p ScratchPeriodList1) Less(i, j int) bool {
	return p[i].Period.Start.Before(p[j].Period.Start)
}
func (p ScratchPeriodList1) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type ScratchPeriod1 struct {
	Period
	ItemList []int   //a list of item user can get at this period
}

func (p ScratchPeriod1) GetPeriod() Period {
	return p.Period
}
func TestGetPeriodFromGenericSortedList(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	periodList1:= ScratchPeriodList1{
		{
			Period:Period{Start:MustFromMysqlFormat("2001-01-00 23:30:00"),End:MustFromMysqlFormat("2001-01-01 23:30:00")},
			ItemList: []int{1,2},
		},
		{
			Period:Period{Start:MustFromMysqlFormat("2001-01-03 23:30:00"),End:MustFromMysqlFormat("2001-01-04 23:30:00")},
			ItemList: []int{2,3},
		},
		{
			Period:Period{Start:MustFromMysqlFormat("2001-01-02 23:30:00"),End:MustFromMysqlFormat("2001-01-03 23:30:00")},
			ItemList: []int{3,4},
		},
	}
	sort.Sort(periodList1)
	i,err:=GetPeriodFromGenericSortedList(MustFromMysqlFormat("2001-01-01 23:00:00"),periodList1)
	t.Equal(err,nil)
	t.Equal(periodList1[i].ItemList,[]int{1,2})
	i,err=GetPeriodFromGenericSortedList(MustFromMysqlFormat("2001-01-03 23:00:00"),periodList1)
	t.Equal(err,nil)
	t.Equal(periodList1[i].ItemList,[]int{3,4})
}
