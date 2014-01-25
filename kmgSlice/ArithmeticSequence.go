package kmgSlice

import (
	"fmt"
)

/*
等差数列
print from start to end ,not include end.
it is valid input ,if step <0
it will panic will invalid input
if start==end,it will output []int{}
*/
func ArithmeticSequence(start int, end int, step int) (output []int) {
	switch {
	case start == end:
		return []int{}
	case step == 0:
		panic("[ArithmeticSequence] step==0")
	case step > 0 && start < end:
		panic(fmt.Errorf("[ArithmeticSequence] start:%d<end:%d step:%d>0", start, end, step))
	case step < 0 && start > end:
		panic(fmt.Errorf("[ArithmeticSequence] start:%d>end:%d step:%d<0", start, end, step))
	}
	for i := start; i < end; i += step {
		output = append(output, i)
	}
	return
}

// a slice with content 0->n-1
func IntRangeSlice(n int) []int {
	if n <= 0 {
		panic(fmt.Errorf("[IntNSlice] n:%d<=0", n))
	}
	output := make([]int, n)
	for i := 0; i < n; i++ {
		output[i] = i
	}
	return output
}
