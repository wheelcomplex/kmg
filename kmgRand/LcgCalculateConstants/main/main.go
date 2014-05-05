package main

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgRand/LcgCalculateConstants"
	"os"
	"strconv"
)

func main() {
	m := uint64(0)
	c := uint64(0)
	if len(os.Args) >= 2 {
		mint, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		m = uint64(mint)
	} else if len(os.Args) == 3 {
		cint, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		c = uint64(cint)
	} else {
		fmt.Printf("usage: %s [m(the range of lcg)] [c]\n", os.Args[0])
		os.Exit(-1)
	}
	LcgCalculateConstants.LcgCalculateConstantsDebug(m, c)
}
