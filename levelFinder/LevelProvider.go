package levelFinder

type LevelProvider interface {
	//到达这个等级lv需要的累计经验
	//默认等级1,需要0累计经验
	//输入范围 [1,MaxLevel()] ,两端包含
	GetExpByLevel(lv int) int
	//最大等级(到达这个等级以后,不能升级了)
	//对某一个实例,这个值不应该变化
	MaxLevel() int
}
