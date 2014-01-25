package kmgRand

import (
//"fmt"
)

func (t *Tester) TestCombinatoricsRandom2d() {
	r := NewInt64SeedKmgRand(0)
	for testcaseId, testcase := range []struct {
		randomer    *CombinatoricsRandom2d
		retLen      int
		retANumList []int
		retBNumList []int
		retHasSolve bool
	}{
		{ //0
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{1, 2},
				BNumList: []int{1, 2},
				ValidCombine: [][]bool{
					[]bool{true, true},
					[]bool{true, true},
				},
			},
			retLen:      3,
			retANumList: []int{1, 2},
			retBNumList: []int{1, 2},
			retHasSolve: true,
		},
		{ //1
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{4, 5},
				ValidCombine: [][]bool{
					[]bool{true, false},
					[]bool{false, true},
					[]bool{true, false},
				},
			},
			retLen:      9,
			retANumList: []int{2, 3, 4},
			retBNumList: []int{4, 5},
			retHasSolve: false,
		},
		{ //2
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{4, 5},
				ValidCombine: [][]bool{
					[]bool{true, true},
					[]bool{false, true},
					[]bool{true, false},
				},
			},
			retLen:      9,
			retANumList: []int{2, 3, 4},
			retBNumList: []int{4, 5},
			retHasSolve: true,
		},
		{ //3
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{3, 5},
				ValidCombine: [][]bool{
					[]bool{true, true},
					[]bool{false, true},
					[]bool{true, false},
				},
			},
			retLen:      8,
			retANumList: []int{2, 3, 3},
			retBNumList: []int{3, 5},
			retHasSolve: true,
		},
		{ //4
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{2, 3, 4},
				BNumList: []int{4, 4},
				ValidCombine: [][]bool{
					[]bool{true, true},
					[]bool{false, true},
					[]bool{true, false},
				},
			},
			retLen:      8,
			retANumList: nil,
			retBNumList: []int{4, 4},
			retHasSolve: true,
		},
		{ //5
			randomer: &CombinatoricsRandom2d{
				ANumList: []int{10, 10, 10},
				BNumList: []int{4, 4},
				ValidCombine: [][]bool{
					[]bool{true, true},
					[]bool{false, true},
					[]bool{true, false},
				},
			},
			retLen:      8,
			retANumList: nil,
			retBNumList: []int{4, 4},
			retHasSolve: true,
		},
	} {
		for i := 0; i < 10; i++ {
			randomer := testcase.randomer
			err := randomer.Random(r)
			if !testcase.retHasSolve {
				t.Ok(err != nil)
				continue
			}
			t.Equal(err, nil)
			t.Equal(len(randomer.Output), testcase.retLen)
			ANumList := make([]int, len(randomer.ANumList))
			BNumList := make([]int, len(randomer.BNumList))
			for _, row := range randomer.Output {
				ANumList[row.X]++
				BNumList[row.Y]++
			}
			//fmt.Println(randomer.Output)
			if testcase.retANumList != nil {
				t.Equal(ANumList, testcase.retANumList)
			}
			if testcase.retBNumList != nil {
				t.EqualMsg(BNumList, testcase.retBNumList, "BNumList not correct testcaseId: %d", testcaseId)
			}
		}
	}
}
