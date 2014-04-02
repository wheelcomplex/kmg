package kmgExcel

import "sort"

type titleArrayFirstSort [][]string

func (ta titleArrayFirstSort) Len() int {
	return len(ta) - 1
}
func (ta titleArrayFirstSort) Less(i, j int) bool {
	return ta[i+1][0] < ta[j+1][0]
}
func (ta titleArrayFirstSort) Swap(i, j int) {
	ta[i+1], ta[j+1] = ta[j+1], ta[i+1]
}

func TitleArraySortByFirst(titleArray [][]string) {
	sort.Sort(titleArrayFirstSort(titleArray))
}
