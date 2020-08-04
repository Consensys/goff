package asm

func generateNeg() {
	fnHeader("Neg", 0, 16)

	// labels
	nonZero := newLabel()

	// registers
	x := popRegister()
	r := popRegister()
	q := popRegister()
	t := popRegisters(nbWords)

	movq("res+0(FP)", r)
	movq("x+8(FP)", x)

	// t = x
	_mov(x, t)

	// x = t[0] | ... | t[n]
	movq(t[0], x)
	for i := 1; i < nbWords; i++ {
		orq(t[i], x)
	}

	testq(x, x)

	// if x != 0, we jump to nonzero label
	jne(nonZero)
	// if x == 0, we set the result to zero and return
	for i := 0; i < nbWords/2; i++ {
		movq(x, r.at(i))
	}
	ret()

	label(nonZero)

	// z = x - q
	for i := 0; i < nbWords; i++ {
		movq(modulus[i], q)
		if i == 0 {
			subq(t[i], q)
		} else {
			sbbq(t[i], q)
		}
		movq(q, r.at(i))
	}

	ret()

}
