package asm

import (
	"fmt"
)

func generateReduce() {
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("reduce", stackSize, 8)

	// registers
	r := popRegister()
	t := popRegisters(nbWords)

	movq("res+0(FP)", r)

	_mov(r, t)
	_reduce(t, r)
	ret()
}

func _reduce(t []register, result interface{}, rOffset ...int) {
	if nbWords > smallModulus {
		reduceLarge(t, result, rOffset...)
		return
	}
	// u = t - q
	u := popRegisters(nbWords)

	_mov(t, u)
	for i := 0; i < nbWords; i++ {
		if i == 0 {
			subq(qAt(i), u[i])
		} else {
			sbbq(qAt(i), u[i])
		}
	}

	// conditional move of u into t (if we have a borrow we need to return t - q)
	for i := 0; i < nbWords; i++ {
		cmovqcc(u[i], t[i])
	}

	// return t
	offset := 0
	if len(rOffset) > 0 {
		offset = rOffset[0]
	}
	_mov(t, result, 0, offset)

	pushRegister(u...)
}

func reduceLarge(t []register, result interface{}, rOffset ...int) {
	// u = t - q
	u := make([]string, nbWords)

	for i := 0; i < nbWords; i++ {
		// use stack
		u[i] = fmt.Sprintf("t%d-%d(SP)", i, 8+i*8)
		movq(t[i], u[i])

		if i == 0 {
			subq(qAt(i), t[i])
		} else {
			sbbq(qAt(i), t[i])
		}
	}

	// conditional move of u into t (if we have a borrow we need to return t - q)
	for i := 0; i < nbWords; i++ {
		cmovqcs(u[i], t[i])
	}

	offset := 0
	if len(rOffset) > 0 {
		offset = rOffset[0]
	}
	// return t
	_mov(t, result, 0, offset)
}
