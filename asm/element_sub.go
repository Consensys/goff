package asm

func generateSub() {
	fnHeader("sub", 0, 24)

	// registers
	t := popRegisters(nbWords)
	x := popRegister()
	y := popRegister()

	movq("x+8(FP)", x)
	_mov(x, t)

	// z = x - y mod q
	movq("y+16(FP)", y)
	_sub(y, t)

	if nbWords > 6 {
		_reduceAfterSub(t, false)
	} else {
		_reduceAfterSub(t, true)
	}

	r := popRegister()
	movq("res+0(FP)", r)
	_mov(t, r)

	ret()

}

func _reduceAfterSub(t []register, noJump bool) {
	if noJump {
		q := popRegisters(nbWords)
		r := popRegister()
		_mov(modulus, q)
		movq(0, r)
		// overwrite with 0 if borrow is set
		for i := 0; i < nbWords; i++ {
			cmovqcc(r, q[i])
		}

		// add registers (q or 0) to t, and set to result
		_add(q, t)

		pushRegister(r)
		pushRegister(q...)
	} else {
		noReduce := newLabel()

		jcc(noReduce)
		for i := 0; i < nbWords; i++ {
			if i == 0 {
				addq(qAt(i), t[i])
			} else {
				adcq(qAt(i), t[i])
			}
		}
		label(noReduce)

	}
}
