package levelFinder

type ArrayLevelProvider []int

func (provider ArrayLevelProvider) GetExpByLevel(lv int) int {
	return provider[lv-1]
}
func (provider ArrayLevelProvider) MaxLevel() int {
	return len(provider) + 1
}
func (provider ArrayLevelProvider) SetExpByLevel(lv int, exp int) {
	provider[lv-1] = exp
}
func NewArrayLevelProvider(maxLevel int) ArrayLevelProvider {
	return make(ArrayLevelProvider, maxLevel-2)
}
