package levelFinder

type LevelProvider interface {
	//根据等级获得升级经验(即累计经验达到这个经验升lv+1级)
	//当前经验为升级经验时,处于lv+1
	//默认升到1级,需要0经验(初始1级,0经验)
	//输入范围 [1,MaxLevel()-1] ,两端包含
	//不处理边界情况,调用者保证边界情况
	GetExpByLevel(lv int) int
	//最大等级(到达这个等级以后,不能升级了)
	//对某一个实例,这个值不应该变化
	MaxLevel() int
}

type LevelExpResult struct {
	Exp                   int  //规整经验(处理负经验问题和满经验问题)
	Level                 int  //等级
	CurrentLevelExcessExp int  // 当前等级溢出经验(满经验时,这个值为0)
	NextLevelAllNeedExp   int  // 从lv 0 exp升到 lv+1 0exp需要的经验(满经验时,这个值为0)
	IsMaxLevel            bool //是否满级
}

//根据经验获取等级,
//合理处理边界情况,不报错
//负累计经验重置为0累计经验
//累计经验超过满级经验重置为 满级经验
func GetLevelByExp(provider LevelProvider, exp int) (result LevelExpResult) {
	if exp < 0 {
		result.Exp = 0
	} else {
		result.Exp = exp
	}
	//找当前等级
	maxLevel := provider.MaxLevel()
	if maxLevel == 1 {
		result.Exp = 0
		result.IsMaxLevel = true
		result.Level = 1
		return
	}
	level := maxLevel
	for i := 1; i <= maxLevel-1; i++ {
		thisLevelExp := provider.GetExpByLevel(i)
		if exp >= thisLevelExp {
			continue
		}
		level = i
		break
	}
	result.Level = level
	if level == maxLevel {
		result.Exp = provider.GetExpByLevel(maxLevel - 1)
		result.IsMaxLevel = true
		return
	}
	if level == 1 {
		result.CurrentLevelExcessExp = result.Exp
		result.NextLevelAllNeedExp = provider.GetExpByLevel(1)
		return
	}
	result.CurrentLevelExcessExp = result.Exp - provider.GetExpByLevel(level-1)
	result.NextLevelAllNeedExp = provider.GetExpByLevel(level) - provider.GetExpByLevel(level-1)
	return
}

//根据等级获取 刚升到这个等级的累计经验
//合理处理边界情况,不报错
//含义和LevelProvider.GetExpByLevel不一样
func GetExpByLevel(provider LevelProvider, level int) (exp int) {
	if level <= 1 {
		return 0
	}

	if level > provider.MaxLevel() {
		level = provider.MaxLevel()
	}
	return provider.GetExpByLevel(level - 1)
}
