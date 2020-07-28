package asm

func generateSub() {
	fnHeader("Sub"+elementName, 0, 24)

	// registers
	t := popRegisters(nbWords)
	x := popRegister()
	y := popRegister()
	r := popRegister()

	movq("x+8(FP)", x)
	_mov(x, t)

	// z = x - y mod q
	movq("y+16(FP)", y)
	_sub(y, t)

	if nbWords > 6 {
		// non constant time, using jumps
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

	} else {
		q := popRegisters(nbWords)
		_mov(modulus, q)
		movq(0, r)
		// overwrite with 0 if borrow is set
		for i := 0; i < nbWords; i++ {
			cmovqcc(r, q[i])
		}

		// add registers (q or 0) to t, and set to result
		_add(q, t)
	}

	movq("res+0(FP)", r)
	_mov(t, r)

	ret()

}

func generateSubE2() {
	fnHeader("Sub"+elementName+"2", 0, 24)

	// registers
	t := popRegisters(nbWords)
	x := popRegister()
	y := popRegister()
	r := popRegister()

	movq("x+8(FP)", x)
	movq("y+16(FP)", y)

	_mov(x, t)

	// set DX to 0
	xorq(r, r)

	// z = x - y mod q
	// move t = x
	_sub(y, t)

	if nbWords > 6 {
		// non constant time, using jumps
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

	} else {
		q := popRegisters(nbWords)
		_mov(modulus, q)
		movq(0, r)
		// overwrite with 0 if borrow is set
		for i := 0; i < nbWords; i++ {
			cmovqcc(r, q[i])
		}

		// add registers (q or 0) to t, and set to result
		_add(q, t)
		pushRegister(q...)
	}

	movq("res+0(FP)", r)

	_mov(t, r)

	_mov(x, t, nbWords)

	// set DX to 0
	xorq(r, r)

	// z = x - y mod q
	// move t = x
	_sub(y, t, nbWords)

	if nbWords > 6 {
		// non constant time, using jumps
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

	} else {
		q := popRegisters(nbWords)
		_mov(modulus, q)
		movq(0, r)
		// overwrite with 0 if borrow is set
		for i := 0; i < nbWords; i++ {
			cmovqcc(r, q[i])
		}

		// add registers (q or 0) to t, and set to result
		_add(q, t)
		pushRegister(q...)
	}

	movq("res+0(FP)", r)

	_mov(t, r, 0, nbWords)

	ret()

}
