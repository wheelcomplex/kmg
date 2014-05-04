package kmgRand

import "fmt"

//TODO 需要高精度计算?
type LcgTransformer struct {
	Start uint64 //start of target range(included)
	Range uint64 //length of target range
	A     uint64 //lcg parameter a
	C     uint64 //lcg parameter c
}

//在范围内进行变换,如果超出规定的范围会报错
//map a num from [0,t.Range-1] to [Start,Start+Range-1] .
func (t LcgTransformer) GenerateInRange(i uint64) (output uint64, err error) {
	if i >= t.Range {
		return 0, fmt.Errorf("[LcgTransformer.GenerateInRange] i[%d]>=t.Range[%d]", i, t.Range)
	}
	if i < 0 {
		return 0, fmt.Errorf("[LcgTransformer.GenerateInRange] i[%d]<0", i)
	}
	return t.Start + (i*t.A+t.C)%t.Range, nil
}
