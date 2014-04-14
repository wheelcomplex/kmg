package kmgTime

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"
)

// a period time start from Start,end to End,
// start must before or equal end
type Period struct {
	Start time.Time
	End   time.Time
}

type PeriodGetter interface {
	GetPeriod() Period
}

var ReflectTypePeriodGetter = reflect.TypeOf((*PeriodGetter)(nil)).Elem()
var NotFoundError = errors.New("not found")

func (p Period) IsIn(t time.Time) bool {
	if t.After(p.End) {
		return false
	}
	if t.Before(p.Start) {
		return false
	}
	return true
}

func (p Period) IsValid() bool {
	if p.End.Before(p.Start) {
		return false
	} else {
		return true
	}
}

func leq(a, b time.Time) bool {
	return a.Before(b) || a.Equal(b)
}

func (p Period) Overlaps(b Period) bool {
	return (leq(b.Start, p.Start) && leq(p.Start, b.End)) ||
		(leq(p.Start, b.Start) && leq(b.Start, p.End))
}

// start must before or equal end
func NewPeriod(Start time.Time, End time.Time) (period Period, err error) {
	if Start.After(End) {
		err = fmt.Errorf("[kmgTime.NewPeriod] Start.After(End) Start:%s End:%s", Start, End)
		return
	}
	return Period{Start: Start, End: End}, nil
}

//SortedList should sort by start time and should not overlap each other
func GetPeriodFromSortedList(t time.Time, SortedList []Period) (index int, ok bool) {
	n := len(SortedList)
	i := sort.Search(n, func(i int) bool {
		return SortedList[i].End.After(t)
	})
	if i == n {
		return 0, false
	}
	if !(SortedList[i].Start.Before(t) || SortedList[i].Start.Equal(t)) {
		return 0, false
	}
	return i, true
}

//实现1,反射
func GetPeriodFromGenericSortedList(t time.Time, SortedList interface{}) (index int, err error) {
	reflectList := reflect.Indirect(reflect.ValueOf(SortedList))
	if reflectList.Kind() != reflect.Slice && reflectList.Kind() != reflect.Array {
		panic(fmt.Errorf("[GetPeriodFromGenericSortedList] need array or slice get %s", reflectList.Kind().String()))
	}
	if !reflectList.Type().Elem().Implements(ReflectTypePeriodGetter) {
		panic(fmt.Errorf("[GetPeriodFromGenericSortedList] need elem implement 'PeriodGetter' get %s",
			reflectList.Elem().Type().Name()))
	}
	n := reflectList.Len()
	i := sort.Search(n, func(i int) bool {
		return reflectList.Index(i).Interface().(PeriodGetter).GetPeriod().End.After(t)
	})
	if i == n {
		return 0, NotFoundError
	}
	if !reflectList.Index(i).Interface().(PeriodGetter).GetPeriod().Start.Before(t) {
		return 0, NotFoundError
	}
	return i, nil
}

type PeriodSlice []Period

func (p PeriodSlice) Len() int {
	return len(p)
}

func (p PeriodSlice) Less(i, j int) bool {
	return p[i].Start.Before(p[j].Start)
}
func (p PeriodSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func PeriodSort(p []Period) {
	sort.Sort(PeriodSlice(p))
}

//实现2,interface
type PeriodListInterface interface {
	Len() int
	GetPeriodAtIndex(i int) Period
	Swap(i, j int)
}

func SelectPeriodFromSortedPeriodList(t time.Time, PeriodList PeriodListInterface) (index int, ok bool) {
	n := PeriodList.Len()
	i := sort.Search(n, func(i int) bool {
		return PeriodList.GetPeriodAtIndex(i).End.After(t)
	})
	if i == n {
		return 0, false
	}
	currentPeriod := PeriodList.GetPeriodAtIndex(i)
	if !(currentPeriod.Start.Before(t) || currentPeriod.Start.Equal(t)) {
		return 0, false
	}
	return i, true
}

type periodListSorter struct {
	PeriodListInterface
}

func (p periodListSorter) Less(i, j int) bool {
	return p.GetPeriodAtIndex(i).Start.Before(p.GetPeriodAtIndex(j).Start)
}
func PeriodListSort(PeriodList PeriodListInterface) {
	sort.Sort(periodListSorter{PeriodList})
}

//实现3,数据结构分离
type PeriodList []PeriodListElement

type PeriodListElement struct {
	Period
	OriginIndex int
}

func (p PeriodList) Len() int {
	return len(p)
}
func (p PeriodList) Less(i, j int) bool {
	return p[i].Period.Start.Before(p[j].Period.Start)
}
func (p PeriodList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PeriodList) Sort() {
	sort.Sort(p)
}

//sort it first
func (p PeriodList) SelectFromTime(t time.Time) (OriginIndex int, ok bool) {
	n := p.Len()
	i := sort.Search(n, func(i int) bool {
		return p[i].End.After(t)
	})
	if i == n {
		return 0, false
	}
	currentPeriod := p[i]
	if !(currentPeriod.Start.Before(t) || currentPeriod.Start.Equal(t)) {
		return 0, false
	}
	return p[i].OriginIndex, true
}

//sort it first
func (p PeriodList) IsIncludeOverlap() bool {
	for i := 0; i < len(p)-1; i++ {
		if p[i].Overlaps(p[i+1].Period) {
			return true
		}
	}
	return false
}
func (p PeriodList) SetAtIndex(period Period, i int) {
	p[i] = PeriodListElement{Period: period, OriginIndex: i}
}
