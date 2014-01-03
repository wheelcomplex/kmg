package levelFinder

import (
	"fmt"
)

type maxLevelLimitProvider struct {
	LevelProvider
	maxLevel int
}

func (provider *maxLevelLimitProvider) MaxLevel() int {
	return provider.maxLevel
}
func MaxLevelLimitProvider(inProvider LevelProvider, MaxLevel int) (outProvider LevelProvider, err error) {
	inMaxLevel := inProvider.MaxLevel()
	if inMaxLevel == MaxLevel {
		return inProvider, nil
	}
	if inMaxLevel < MaxLevel {
		return nil, fmt.Errorf("[MaxLevelLimitProvider] inProvider.MaxLevel() %d <MaxLevel %d", inMaxLevel, MaxLevel)
	}
	if MaxLevel <= 0 {
		return nil, fmt.Errorf("[MaxLevelLimitProvider] MaxLevel %d <= 0", MaxLevel)
	}
	return &maxLevelLimitProvider{LevelProvider: inProvider, maxLevel: MaxLevel}, nil
}
