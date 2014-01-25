package kmgSort

import "sort"

type IntLessCallbackSortT struct {
	Data     []int
	LessFunc func(a int, b int) bool
}

func (s *IntLessCallbackSortT) Len() int {
	return len(s.Data)
}

func (s *IntLessCallbackSortT) Less(i, j int) bool {
	return s.LessFunc(i, j)
}

func (s *IntLessCallbackSortT) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

func (s *IntLessCallbackSortT) Sort() {
	sort.Sort(s)
}

func IntLessCallbackSort(Data []int, LessFunc func(a int, b int) bool) {
	sort.Sort(&IntLessCallbackSortT{Data: Data, LessFunc: LessFunc})
}
