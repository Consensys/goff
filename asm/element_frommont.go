package asm

func generateFromMont() {
	stackSize := 8
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("FromMont", stackSize, 8, dx, ax)
	writeLn("NO_LOCAL_POINTERS")
	writeLn(`
	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// when y = 1 we have: 
	// for i=0 to N-1
	// 		t[i] = x[i]
	// for i=0 to N-1
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 		    (C,t[j-1]) := t[j] + m*q[j] + C
	// 		t[N-1] = C`)

	noAdx := newLabel()
	// check ADX instruction support
	cmpb("·supportAdx(SB)", 1)
	jne(noAdx)

	// registers
	t := popRegisters(nbWords)
	r := popRegister()

	movq("res+0(FP)", r)

	// 	for i=0 to N-1
	//     t[i] = a[i]
	_mov(r, t)

	var tmp register
	hasRegisters := availableRegisters() > 0
	if !hasRegisters {
		tmp = r
	} else {
		tmp = popRegister()
	}
	for i := 0; i < nbWords; i++ {

		xorq(dx, dx)

		// m := t[0]*q'[0] mod W
		regM := dx
		movq(t[0], dx)
		mulxq(qInv0(), regM, ax, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		xorq(ax, ax)

		// C,_ := t[0] + m*q[0]
		comment("C,_ := t[0] + m*q[0]")

		mulxq(qAt(0), ax, tmp)
		adcxq(t[0], ax)
		movq(tmp, t[0])

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
		adoxq(ax, t[nbWordsLastIndex])

	}

	if !hasRegisters {
		movq("res+0(FP)", r)
	} else {
		pushRegister(tmp)
	}
	// ---------------------------------------------------------------------------------------------
	// reduce
	_reduce(t, r)
	ret()

	// No adx
	label(noAdx)
	movq("res+0(FP)", ax)
	movq(ax, "(SP)")
	writeLn("CALL ·_fromMontGeneric(SB)")
	ret()

}
