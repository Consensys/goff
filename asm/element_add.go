package asm

func generateAdd() {
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("Add", stackSize, 24)

	// registers
	x := popRegister()
	y := popRegister()
	r := popRegister()
	t := popRegisters(nbWords)

	movq("x+8(FP)", x)

	// t = x
	_mov(x, t)

	movq("y+16(FP)", y)

	// t = t + y = x + y
	_add(y, t)

	// dereference res
	movq("res+0(FP)", r)

	// reduce t into res
	_reduce(t, r)

	ret()

}

func generateAddE2() {
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("Add"+"2", stackSize, 24)

	// registers
	x := popRegister()
	y := popRegister()
	r := popRegister()
	t := popRegisters(nbWords)

	movq("x+8(FP)", x)

	// move t = x
	_mov(x, t)

	movq("y+16(FP)", y)

	// t = t + y = x + y
	_add(y, t)

	// reduce
	movq("res+0(FP)", r)
	_reduce(t, r)

	// move x+offset(nbWords) into t
	_mov(x, t, nbWords)

	// add y+offset(nbWords) into t
	_add(y, t, nbWords)

	// reduce t into r with offset nbWords
	_reduce(t, r, nbWords)

	ret()

}
