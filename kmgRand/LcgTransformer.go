package kmgRand

import "fmt"

type LcgTransformer struct {
	Start uint64 //start of target range(included)
	Range uint64 //length of target range
	A     uint64 //lcg parameter a
	C     uint64 //lcg parameter c
}

//在范围内进行变换,如果超出规定的范围会报错
//map a num from [0,t.Range-1] to [Start,Start+Range-1] .
func (t LcgTransformer) GenerateInRange(i uint64) (output uint64) {
	if i >= t.Range {
		panic(fmt.Errorf("[LcgTransformer.GenerateInRange] i[%d]>=t.Range[%d]", i, t.Range))
	}
	return t.Start + (i*t.A+t.C)%t.Range
}

func (t LcgTransformer) Generate(i uint64) (output uint64) {
	return t.Start + (i*t.A+t.C)%t.Range
}
