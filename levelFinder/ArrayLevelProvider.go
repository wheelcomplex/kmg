package levelFinder

import "fmt"

//使用数组储存 等级-经验 信息
// key+2 表示升级到该等级需要的累计经验
type ArrayLevelProvider []int

func (provider ArrayLevelProvider) GetExpByLevel(lv int) int {
	if lv > provider.MaxLevel() || lv <= 0 {
		panic(fmt.Errorf("[ArrayLevelProvider.GetExpByLevel] lv[%d]>provider.MaxLevel()[%d] || lv[%d]<=0", lv, provider.MaxLevel(), lv))
	}
	if lv == 1 {
		return 0
	}
	return provider[lv-2]
}
func (provider ArrayLevelProvider) MaxLevel() int {
	return len(provider) + 1
}
func (provider ArrayLevelProvider) SetExpByLevel(lv int, exp int) {
	if lv == 1 {
		if exp != 0 {
			panic(fmt.Errorf("[ArrayLevelProvider.SetExpByLevel] lv==1&&exp[%d]!=0", exp))
		}
		return
	}
	provider[lv-2] = exp
}

//maxLevel表示需要表示的最大等级
//maxLevel需要>=1
func NewArrayLevelProvider(maxLevel int) ArrayLevelProvider {
	return make(ArrayLevelProvider, maxLevel-1)
}
