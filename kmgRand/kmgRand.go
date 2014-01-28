package kmgRand

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"sort"
)

func NewCryptSeedMathRand() (r *mrand.Rand, err error) {
	var seed int64
	err = binary.Read(crand.Reader, binary.LittleEndian, &seed)
	if err != nil {
		return
	}
	return mrand.New(mrand.NewSource(seed)), nil
}

//a kmgRand new from crypt source,only use 8 Byte crypt random...
func NewCryptSeedKmgRand() (r *KmgRand, err error) {
	mr, err := NewCryptSeedMathRand()
	if err != nil {
		return
	}
	return &KmgRand{mr}, nil
}

func NewInt64SeedKmgRand(seed int64) (r *KmgRand) {
	mr := mrand.New(mrand.NewSource(seed))
	return &KmgRand{mr}
}

type KmgRand struct {
	*mrand.Rand
}

//return not repeat size number in [0,n) as random order,
//it will panic if size>n or size<0 or n<0
func (r *KmgRand) MulitChoice(totalLength int, choiceNumber int) []int {
	if choiceNumber > totalLength || totalLength < 0 || choiceNumber < 0 {
		panic(fmt.Errorf("[KmgRand.MulitChoice] input error totalLength=%d choiceNumber=%d,need (choiceNumber<=totalLength&&totalLength>=0&&choiceNumber>=0) ",
			totalLength, choiceNumber))
	}
	return r.Perm(totalLength)[:choiceNumber]
}

//return not repeat size number in [0,n) as origin order,
//it will panic if size>n or size<0 or n<0
func (r *KmgRand) MulitChoiceOriginOrder(n int, size int) []int {
	perm := r.MulitChoice(n, size)
	sort.IntSlice(perm).Sort()
	return perm
}

// return true if that event happend in that possibility.
// possibility should be float in [0,1]
func (r *KmgRand) HappendBaseOnPossibility(possibility float64) bool {
	if possibility > 1+1e-10 || possibility < -1e-10 {
		panic(fmt.Errorf("[KmgRand.HappendBaseOnPossibility] possibility:%f > 1 or < 0", possibility))
	}
	out := r.Float64()
	return out < possibility
}

//return a random int in [min,max]
func (r *KmgRand) IntBetween(min int, max int) int {
	if min > max {
		panic(fmt.Errorf("[KmgRand.IntBetween] min:%d<max:%d", min, max))
	} else if min == max {
		return min
	}
	return r.Intn(max-min) + min
}

func (r *KmgRand) ChoiceFromIntSlice(slice []int) int {
	return slice[r.Intn(len(slice))]
}

func (r *KmgRand) PermIntSlice(slice []int) (output []int) {
	thisLen := len(slice)
	output = make([]int, thisLen)
	permSlice := r.Perm(thisLen)
	for i := 0; i < thisLen; i++ {
		output[i] = slice[permSlice[i]]
	}
	return
}
