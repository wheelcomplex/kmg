package kmgSlice

func (t *Tester) TestIntSliceRemoveAt() {
	s := []int{1, 2, 3}
	IntSliceRemoveAt(&s, 1)
	t.Equal(s, []int{1, 3})
}
