package levelFinder

import (
	"sort"
)

type LevelExpResult struct {
	Exp                   int  //规整经验(已处理负经验问题和满经验问题,请返回后重新写回当前经验)
	Level                 int  //等级
	CurrentLevelExcessExp int  // 当前等级溢出经验(满级时,这个值为0)
	NextLevelAllNeedExp   int  // 从lv 0 exp升到 lv+1 0exp需要的经验(满级时,这个值为0)
	IsMaxLevel            bool //是否满级
}

//根据经验获取等级,
//合理处理边界情况,不报错
//负累计经验重置为0
//累计经验超过满级时,将累计经验重置为 满级经验
func GetLevelByExp(provider LevelProvider, exp int) (result LevelExpResult) {
	if exp < 0 {
		exp = 0
	}
	result.Exp = exp
	//找当前等级
	maxLevel := provider.MaxLevel()
	if maxLevel == 1 {
		result.Exp = 0
		result.IsMaxLevel = true
		result.Level = 1
		return
	}
	level := searchLevelFromProvider(provider, exp)
	result.Level = level
	if level == maxLevel {
		result.Exp = provider.GetExpByLevel(maxLevel)
		result.IsMaxLevel = true
		return
	}
	if level == 1 {
		result.CurrentLevelExcessExp = result.Exp
		result.NextLevelAllNeedExp = provider.GetExpByLevel(2)
		return
	}
	result.CurrentLevelExcessExp = result.Exp - provider.GetExpByLevel(level)
	result.NextLevelAllNeedExp = provider.GetExpByLevel(level+1) - provider.GetExpByLevel(level)
	return
}

/*
	从经验找等级
*/
func searchLevelFromProvider(provider LevelProvider, exp int) (level int) {
	maxLevel := provider.MaxLevel()
	//i+1 for start from 1    TODO what's this stuff mean?
	level = sort.Search(maxLevel-1, func(i int) bool {
		return exp < provider.GetExpByLevel(i+1)
	})
	if level == maxLevel-1 && exp >= provider.GetExpByLevel(maxLevel) {
		return maxLevel
	}
	return level
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
	return provider.GetExpByLevel(level)
}
