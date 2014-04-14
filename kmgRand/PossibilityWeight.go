package kmgRand

import "sort"

//please use NewPossibilityWeightRander, field expose only for serialize
type PossibilityWeightRander struct {
	SearchList []int
	Total      int
}

func (p PossibilityWeightRander) ChoiceOne(r *KmgRand) (Index int) {
	intR := r.IntBetween(0, p.Total-1)
	return sort.SearchInts(p.SearchList, intR)
}

func NewPossibilityWeightRander(weightList []int) PossibilityWeightRander {
	pwl := PossibilityWeightRander{
		SearchList: make([]int, len(weightList)),
	}
	for i, weight := range weightList {
		pwl.SearchList[i] = pwl.Total
		pwl.Total += weight
	}
	return pwl
}
