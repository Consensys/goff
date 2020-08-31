package asm

func mulAdx(yat, xat func(int) string, uglyHook func(int)) []register {
	// registers
	t := popRegisters(nbWords)
	A := popRegister()

	for i := 0; i < nbWords; i++ {
		xorq(dx, dx)

		movq(yat(i), dx)
		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < nbWords; j++ {
			xj := xat(j)

			reg := A
			if i == 0 {
				if j == 0 {
					mulxq(xj, t[j], t[j+1])
				} else if j != nbWordsLastIndex {
					reg = t[j+1]
				}
			} else if j != 0 {
				adcxq(A, t[j])
			}

			if !(i == 0 && j == 0) {
				mulxq(xj, ax, reg)
				adoxq(ax, t[j])
			}
		}
		if uglyHook != nil {
			uglyHook(i)
		}

		comment("add the last carries to " + string(A))
		movq(0, dx)
		adcxq(dx, A)
		adoxq(dx, A)

		// m := t[0]*q'[0] mod W
		m := dx
		movq(t[0], dx)
		mulxq(qInv0(), m, ax, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		xorq(ax, ax)

		// C,_ := t[0] + m*q[0]
		comment("C,_ := t[0] + m*q[0]")

		needPop := false
		if availableRegisters() == 0 {
			needPop = true
			pushq(A)
			pushRegister(A)
		}
		tmp := popRegister()
		mulxq(qAt(0), ax, tmp)
		adcxq(t[0], ax)
		movq(tmp, t[0])
		pushRegister(tmp)
		if needPop {
			A = popRegister()
			popq(A)
		}

		comment("for j=1 to N-1")
		comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < nbWords; j++ {
			adcxq(t[j], t[j-1])
			mulxq(qAt(j), ax, t[j])
			adoxq(ax, t[j-1])
		}
		movq(0, ax)
		adcxq(ax, t[nbWordsLastIndex])
		adoxq(A, t[nbWordsLastIndex])
	}

	// free registers
	pushRegister(A)

	return t
}

func generateInnerMul(isSquare bool) {

	noAdx := newLabel()

	// check ADX instruction support
	cmpb("·supportAdx(SB)", 1)
	jne(noAdx)
	{
		var t []register
		if isSquare {
			x := popRegister()
			movq("x+8(FP)", x)

			xat := func(i int) string {
				return x.at(i)
			}
			t = mulAdx(xat, xat, nil)
			pushRegister(x, x)
		} else {
			x := popRegister()
			y := popRegister()
			movq("x+8(FP)", x)
			movq("y+16(FP)", y)

			yat := func(i int) string {
				return y.at(i)
			}
			xat := func(i int) string {
				return x.at(i)
			}
			t = mulAdx(yat, xat, nil)
			pushRegister(x, y)
		}

		r := popRegister()
		// ---------------------------------------------------------------------------------------------
		// reduce
		movq("res+0(FP)", r)
		_reduce(t, r)
		ret()
		pushRegister(r)
	}

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		label(noAdx)
		builder.reset()
		builder.remove(ax)
		builder.remove(dx)
		x := popRegister()
		y := popRegister()
		if isSquare {
			movq("x+8(FP)", x)
			movq("x+8(FP)", y)
		} else {
			movq("x+8(FP)", x)
			movq("y+16(FP)", y)
		}

		mulNoAdx(x, y)
	}
}

func generateMul() {
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("mul", stackSize, 24, dx, ax)
	writeLn(`
	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// however, to benefit from the ADCX and ADOX carry chains
	// we split the inner loops in 2:
	// for i=0 to N-1
	// 		for j=0 to N-1
	// 		    (A,t[j])  := t[j] + x[j]*y[i] + A
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 		    (C,t[j-1]) := t[j] + m*q[j] + C
	// 		t[N-1] = C + A
	`)
	if nbWords > smallModulus {
		generateInnerMulLarge(false)
	} else {
		generateInnerMul(false)
	}

}

func generateInnerMulLarge(isSquare bool) {
	writeLn("NO_LOCAL_POINTERS")
	noAdx := newLabel()
	// check ADX instruction support
	cmpb("·supportAdx(SB)", 1)
	jne(noAdx)

	// registers
	t := popRegisters(nbWords)
	A := popRegister()

	for i := 0; i < nbWords; i++ {

		xorq(dx, dx)
		yi := dx
		if isSquare {
			movq("x+8(FP)", yi)
		} else {
			movq("y+16(FP)", yi)
		}
		movq(yi.at(i), yi)
		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < nbWords; j++ {
			xj := ax
			movq("x+8(FP)", xj)
			movq(xj.at(j), xj)

			reg := A
			if i == 0 {
				if j == 0 {
					mulxq(xj, t[j], t[j+1])
				} else if j != nbWordsLastIndex {
					reg = t[j+1]
				}
			} else if j != 0 {
				adcxq(A, t[j])
			}

			if !(i == 0 && j == 0) {
				mulxq(xj, ax, reg)
				adoxq(ax, t[j])
			}
		}

		comment("add the last carries to " + string(A))
		movq(0, dx)
		adcxq(dx, A)
		adoxq(dx, A)
		pushq(A)

		// m := t[0]*q'[0] mod W
		regM := dx
		movq(t[0], dx)
		mulxq(qInv0(), regM, ax, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		xorq(ax, ax)

		// C,_ := t[0] + m*q[0]
		comment("C,_ := t[0] + m*q[0]")
		mulxq(qAt(0), ax, A)
		adcxq(t[0], ax)
		movq(A, t[0])

		comment("for j=1 to N-1")
		comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < nbWords; j++ {
			adcxq(t[j], t[j-1])
			mulxq(qAt(j), ax, t[j])
			adoxq(ax, t[j-1])
		}

		popq(A)
		movq(0, ax)
		adcxq(ax, t[nbWordsLastIndex])
		adoxq(A, t[nbWordsLastIndex])
	}

	// free registers
	pushRegister(A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := popRegister()
	movq("res+0(FP)", r)
	reduceLarge(t, r)
	ret()

	// No adx
	label(noAdx)
	movq("res+0(FP)", ax)
	movq(ax, "(SP)")
	movq("x+8(FP)", ax)
	movq(ax, "8(SP)")
	if isSquare {
		writeLn("CALL ·_squareGeneric(SB)")
		ret()
	} else {
		movq("y+16(FP)", ax)
		movq(ax, "16(SP)")
		writeLn("CALL ·_mulGeneric(SB)")
		ret()
	}

}

func mulNoAdx(x, y register) {
	// registers
	t := popRegisters(nbWords)
	C := popRegister()
	yi := popRegister()
	A := popRegister()
	m := popRegister()

	for i := 0; i < nbWords; i++ {
		// (A,t[0]) := t[0] + x[0]*y[{{$i}}]
		movq(x.at(0), ax)
		movq(y.at(i), yi)
		mulq(yi)
		if i != 0 {
			addq(ax, t[0])
			adcq(0, dx)
		} else {
			movq(ax, t[0])
		}
		movq(dx, A)

		// m := t[0]*q'[0] mod W
		movq(qInv0(), m)
		imulq(t[0], m)

		// C,_ := t[0] + m*q[0]
		movq(modulus[0], ax)
		mulq(m)
		addq(t[0], ax)
		adcq(0, dx)
		movq(dx, C)

		// for j=1 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < nbWords; j++ {
			movq(x.at(j), ax)
			mulq(yi)
			if i != 0 {
				addq(A, t[j])
				adcq(0, dx)
				addq(ax, t[j])
				adcq(0, dx)
			} else {
				movq(A, t[j])
				addq(ax, t[j])
				adcq(0, dx)
			}
			movq(dx, A)

			movq(modulus[j], ax)
			mulq(m)
			addq(t[j], C)
			adcq(0, dx)
			addq(ax, C)
			adcq(0, dx)
			movq(C, t[j-1])
			movq(dx, C)
		}

		addq(C, A)
		movq(A, t[nbWordsLastIndex])

	}

	// ---------------------------------------------------------------------------------------------
	// reduce
	pushRegister(C, yi, A, m, y)

	movq("res+0(FP)", x)
	_reduce(t, x)
	ret()
}
