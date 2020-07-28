package asm

func generateSquare() {
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("Square"+elementName, stackSize, 16, dx, ax)
	writeLn(`
	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// for i=0 to N-1
	// A, t[i] = x[i] * x[i] + t[i]
	// p = 0
	// for j=i+1 to N-1
	//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
	// m = t[0] * q'[0]
	// C, _ = t[0] + q[0]*m
	// for j=1 to N-1
	//     C, t[j-1] = q[j]*m +  t[j] + C
	// t[N-1] = C + A

	`)
	if nbWords > 6 {
		generateInnerMulLarge(true)
		return
	}
	if !noCarrySquare {
		generateInnerMul(true)
		return
	}

	noAdx := newLabel()
	// check ADX instruction support
	cmpb("·supportAdx(SB)", 1)
	jne(noAdx)

	// registers
	t := popRegisters(nbWords)
	x := popRegister()
	A := popRegister()

	movq("x+8(FP)", x)
	for i := 0; i < nbWords; i++ {

		xorq(ax, ax)

		movq(x.at(i), dx)

		// instead of
		// for j=i+1 to N-1
		//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
		// we first add the x[j]*x[i] to a temporary u (set of registers)
		// set double it, before doing
		// for j=i+1 to N-1
		//     A,t[j] = u[j] + t[j] + A
		if i != nbWordsLastIndex {
			u := make([]register, (nbWords - i - 1))
			for i := 0; i < len(u); i++ {
				u[i] = popRegister()
			}
			offset := i + 1

			// 1- compute u = x[j] * x[i]
			// for j=i+1 to N-1
			//     A,u[j] = x[j]*x[i] + A
			if (i + 1) == nbWordsLastIndex {
				mulxq(x.at(i+1), u[0], A)
			} else {
				for j := i + 1; j < nbWords; j++ {
					yj := x.at(j)
					if j == i+1 {
						// first iteration
						mulxq(yj, u[j-offset], u[j+1-offset])
					} else {
						if j == nbWordsLastIndex {
							mulxq(yj, ax, A)
						} else {
							mulxq(yj, ax, u[j+1-offset])
						}
						adcxq(ax, u[j-offset])
					}
				}
				movq(0, ax)
				adcxq(ax, A)
				xorq(ax, ax)
			}

			if i == 0 {
				// C, t[i] = x[i] * x[i] + t[i]
				mulxq(dx, t[i], dx)

				// when i == 0, T is not set yet
				// so  we can use ADOXQ carry chain to propagate C from x[i] * x[i] + t[i] (dx)

				// for j=i+1 to N-1
				// 		C, t[j] = u[j] + u[j] + t[j] + C
				for j := 0; j < len(u); j++ {
					adcxq(u[j], u[j])
					movq(u[j], t[j+offset])
					if j == 0 {
						adoxq(dx, t[j+offset])
					} else {
						adoxq(ax, t[j+offset])
					}
				}

				adcxq(A, A)
				adoxq(ax, A)

			} else {
				// i != 0 so T is set.
				// we first use ADOXQ carry chain to perform t = u + u + t
				for j := 0; j < len(u); j++ {
					adcxq(u[j], u[j])
					adoxq(u[j], t[j+offset])
				}

				adcxq(A, A)
				adoxq(ax, A)

				// reset flags
				xorq(ax, ax)

				// C, t[i] = x[i] * x[i] + t[i]
				mulxq(dx, ax, dx)
				adoxq(ax, t[i])
				movq(0, ax)

				// propagate C
				for j := i + 1; j < nbWords; j++ {
					if j == i+1 {
						adoxq(dx, t[j])
					} else {
						adoxq(ax, t[j])
					}
				}

				adoxq(ax, A)
			}

			pushRegister(u...)

		} else {
			// i == last index
			mulxq(dx, ax, A)
			adcxq(ax, t[i])
			movq(0, ax)
			adcxq(ax, A)
		}

		tmp := popRegister()
		// m := t[0]*q'[0] mod W
		regM := dx
		movq(t[0], dx)
		mulxq(qInv0(), regM, ax, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		xorq(ax, ax)

		// C,_ := t[0] + m*q[0]
		mulxq(qAt(0), ax, tmp)
		adcxq(t[0], ax)
		movq(tmp, t[0])

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

		pushRegister(tmp)
	}

	// free registers
	pushRegister(x, A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := popRegister()
	movq("res+0(FP)", r)
	_reduce(t, r)
	ret()

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		label(noAdx)
		builder.reset()
		builder.remove(ax)
		builder.remove(dx)
		x := popRegister()
		y := popRegister()
		movq("x+8(FP)", x)
		movq("x+8(FP)", y)
		mulNoAdx(x, y)
	}

}
