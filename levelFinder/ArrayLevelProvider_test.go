package levelFinder

func (t *Tester) TestArrayLevelProvider() {
	arrayLevelProvider := NewArrayLevelProvider(3)
	arrayLevelProvider.SetExpByLevel(1, 0)
	arrayLevelProvider.SetExpByLevel(2, 100)
	arrayLevelProvider.SetExpByLevel(3, 200)
	t.Equal([]int(arrayLevelProvider), []int{100, 200})
	t.Equal(arrayLevelProvider.MaxLevel(), 3)

	arrayLevelProvider = NewArrayLevelProvider(2)
	arrayLevelProvider.SetExpByLevel(1, 0)
	arrayLevelProvider.SetExpByLevel(2, 100)
	t.Equal([]int(arrayLevelProvider), []int{100})
	t.Equal(arrayLevelProvider.MaxLevel(), 2)

	t.Equal(arrayLevelProvider.GetExpByLevel(2), 100)
	t.Equal(arrayLevelProvider.GetExpByLevel(1), 0)

	arrayLevelProvider = NewArrayLevelProvider(1)
	t.Equal(arrayLevelProvider.MaxLevel(), 1)
	t.Equal([]int(arrayLevelProvider), []int{})
}
