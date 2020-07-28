package asm

func generateDouble() {
	// func header
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("Double", stackSize, 16)

	// registers
	x := popRegister()
	r := popRegister()
	t := popRegisters(nbWords)

	movq("res+0(FP)", r)
	movq("x+8(FP)", x)

	_mov(x, t)
	_add(t, t)
	_reduce(t, r)

	ret()
}

func generateDoubleE2() {
	// func header
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("Double"+"2", stackSize, 16)

	// registers
	x := popRegister()
	r := popRegister()
	t := popRegisters(nbWords)

	movq("res+0(FP)", r)
	movq("x+8(FP)", x)

	_mov(x, t)
	_add(t, t)
	_reduce(t, r)
	_mov(x, t, nbWords)
	_add(t, t)
	_reduce(t, r, nbWords)

	ret()
}
