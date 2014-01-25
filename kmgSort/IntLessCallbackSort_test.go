package kmgSort

func (t *Tester) TestIntLessCallbackSort() {
	a := []int{2, 1, 3, 0}
	IntLessCallbackSort(a, func(i int, j int) bool {
		return a[i] < a[j]
	})
	t.Equal(a, []int{0, 1, 2, 3})
}
