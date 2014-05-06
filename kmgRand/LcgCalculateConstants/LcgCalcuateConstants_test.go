package LcgCalculateConstants

import (
	"github.com/bronze1man/kmg/kmgRand"
	"testing"
)

func TestLcgCalculateConstants(ot *testing.T) {
	r := uint64(1e4)
	a, c := LcgCalculateConstants(r, 0)
	f := kmgRand.LcgTransformer{
		Start: 0,
		Range: r,
		A:     a,
		C:     c,
	}
	hasValueMap := make([]bool, r)
	for i := uint64(0); i < r; i++ {
		out := f.GenerateInRange(i)
		hasValueMap[int(out)] = true
	}
	for i := uint64(0); i < r; i++ {
		if hasValueMap[i] != true {
			ot.Fatalf("[BenchmarkLcgCalculateConstants] value %d not generate", i)
		}
	}
	/*
		验证某个实际参数是正确的
		r = 1e8
		f = kmgRand.LcgTransformer{
			Start: 0,
			Range: 1e8,
			A:     671088641,
			C:     2531011,
		}
		hasValueMap = make([]bool, r)
		for i := uint64(0); i < r; i++ {
			out:= f.GenerateInRange(i)
			hasValueMap[int(out)] = true
		}
		for i := uint64(0); i < r; i++ {
			if hasValueMap[i]!=true{
				ot.Fatalf("[BenchmarkLcgCalculateConstants] value %d not generate",i)
			}
		}
	*/
	/*
		   修正计算参数
			r=  uint64(1e8)
			a, c = LcgCalculateConstants(r, 0)
			for ;a<1e9;a=(a-1)*2+1{
				f = kmgRand.LcgTransformer{
					Start: 1e8,
					Range: r,
					A:     a,
					C:     c,
				}
				if !isRandomNess(f){
					continue
				}
				fmt.Println(r,a,c)
				for i := uint64(0); i < 100; i++ {
					out:= f.GenerateInRange(i)
					fmt.Printf("%d ",out)
				}
				fmt.Println()
			}
	*/
}

//需要33.2ns,考虑到仅计算Generate需要20ns,应该没有优化空间了
func BenchmarkLcgCalculateConstants(b *testing.B) {
	//t := kmgTest.NewTestTools(b)
	r := uint64(b.N)
	a, c := LcgCalculateConstants(r, 0)
	f := kmgRand.LcgTransformer{
		Start: 0,
		Range: r,
		A:     a,
		C:     c,
	}
	hasValueMap := make([]bool, r)
	for i := uint64(0); i < r; i++ {
		out := f.GenerateInRange(i)
		hasValueMap[int(out)] = true
	}
	for i := uint64(0); i < r; i++ {
		//t.Equal(hasValueMap[i],true)    //使用快捷方法需要120ns
		if !hasValueMap[i] {
			b.Fatalf("[BenchmarkLcgCalculateConstants] value %d not generate", i)
		}
	}
}
func int64abs(i uint64) uint64 {
	if int64(i) > 0 {
		return i
	} else {
		return uint64(-int64(i))
	}
}
func isRandomNess(f kmgRand.LcgTransformer) bool {
	for _, i := range []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		if int64abs(f.GenerateInRange(i)-f.GenerateInRange(0)) < 1e6 &&
			int64abs(f.GenerateInRange(2*i)-f.GenerateInRange(i)) < 1e6 {
			return false
		}
	}

	return true
}
