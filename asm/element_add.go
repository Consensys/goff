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
