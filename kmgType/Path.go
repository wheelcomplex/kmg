package kmgType

import (
	"strings"
)

//a path to a type in whole type system
type Path []string

func (p Path) String() string {
	return strings.Join(p, ",")
}
func ParsePath(ps string) Path {
	psa := strings.Split(ps, ",")
	pso := []string{}
	for _, v := range psa {
		v = strings.TrimSpace(v)
		pso = append(pso, v)
	}
	return Path(pso)
}
