package LcgCalculateConstants

/*
//cgo not work in cygwin...,and main in c not work in cgo too..
#cgo LDFLAGS: -lm
#include "c/rand-lcg.h"
#include "c/rand-lcg.c"
#include "c/rand-primegen.h"
#include "c/rand-primegen.c"
*/
import "C"

func LcgCalculateConstants(m uint64, inC uint64) (a uint64, c uint64) {
	C.lcg_calculate_constants(C.uint64_t(m), (*C.uint64_t)(&a), (*C.uint64_t)(&inC), 0)
	c = inC
	return
}
