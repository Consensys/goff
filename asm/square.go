package asm

import (
	"github.com/consensys/bavard"
)

func (b *Builder) square(asm *bavard.Assembly) error {
	asm.FuncHeader("_squareADX"+b.elementName, 0, 16)
	asm.WriteLn(`
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

	// check ADX instruction support
	asm.CMPB("·supportAdx(SB)", 1)
	asm.JNE("no_adx")

	// registers
	t := asm.PopRegisters(b.nbWords)
	x := asm.PopRegister()
	A := asm.PopRegister()

	// dereference x
	asm.MOVQ("x+8(FP)", x)
	for i := 0; i < b.nbWords; i++ {

		asm.XORQ(bavard.AX, bavard.AX)

		asm.MOVQ(x.At(i), bavard.DX)

		// instead of
		// for j=i+1 to N-1
		//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
		// we first add the x[j]*x[i] to a temporary u (set of registers)
		// set double it, before doing
		// for j=i+1 to N-1
		//     A,t[j] = u[j] + t[j] + A
		if i != b.nbWordsLastIndex {
			u := make([]bavard.Register, (b.nbWords - i - 1))
			for i := 0; i < len(u); i++ {
				u[i] = asm.PopRegister()
			}
			offset := i + 1

			// 1- compute u = x[j] * x[i]
			// for j=i+1 to N-1
			//     A,u[j] = x[j]*x[i] + A
			if (i + 1) == b.nbWordsLastIndex {
				asm.MULXQ(x.At(i+1), u[0], A)
			} else {
				for j := i + 1; j < b.nbWords; j++ {
					yj := x.At(j)
					if j == i+1 {
						// first iteration
						asm.MULXQ(yj, u[j-offset], u[j+1-offset])
					} else {
						if j == b.nbWordsLastIndex {
							asm.MULXQ(yj, bavard.AX, A)
						} else {
							asm.MULXQ(yj, bavard.AX, u[j+1-offset])
						}
						asm.ADCXQ(bavard.AX, u[j-offset])
					}
				}
				asm.MOVQ(0, bavard.AX)
				asm.ADCXQ(bavard.AX, A)
				asm.XORQ(bavard.AX, bavard.AX)
			}

			if i == 0 {
				// C, t[i] = x[i] * x[i] + t[i]
				asm.MULXQ(bavard.DX, t[i], bavard.DX)

				// when i == 0, T is not set yet
				// so  we can use ADOXQ carry chain to propagate C from x[i] * x[i] + t[i] (dx)

				// for j=i+1 to N-1
				// 		C, t[j] = u[j] + u[j] + t[j] + C
				for j := 0; j < len(u); j++ {
					asm.ADCXQ(u[j], u[j])
					asm.MOVQ(u[j], t[j+offset])
					if j == 0 {
						asm.ADOXQ(bavard.DX, t[j+offset])
					} else {
						asm.ADOXQ(bavard.AX, t[j+offset])
					}
				}

				asm.ADCXQ(A, A)
				asm.ADOXQ(bavard.AX, A)

			} else {
				// i != 0 so T is set.
				// we first use ADOXQ carry chain to perform t = u + u + t
				for j := 0; j < len(u); j++ {
					asm.ADCXQ(u[j], u[j])
					asm.ADOXQ(u[j], t[j+offset])
				}

				asm.ADCXQ(A, A)
				asm.ADOXQ(bavard.AX, A)

				// reset flags
				asm.XORQ(bavard.AX, bavard.AX)

				// C, t[i] = x[i] * x[i] + t[i]
				asm.MULXQ(bavard.DX, bavard.AX, bavard.DX)
				asm.ADOXQ(bavard.AX, t[i])
				asm.MOVQ(0, bavard.AX)

				// propagate C
				for j := i + 1; j < b.nbWords; j++ {
					if j == i+1 {
						asm.ADOXQ(bavard.DX, t[j])
					} else {
						asm.ADOXQ(bavard.AX, t[j])
					}
				}

				asm.ADOXQ(bavard.AX, A)
			}

			asm.PushRegister(u...)

		} else {
			// i == last index
			asm.MULXQ(bavard.DX, bavard.AX, A)
			asm.ADCXQ(bavard.AX, t[i])
			asm.MOVQ(0, bavard.AX)
			asm.ADCXQ(bavard.AX, A)
		}

		tmp := asm.PopRegister()
		// m := t[0]*q'[0] mod W
		regM := bavard.DX
		asm.MOVQ(t[0], bavard.DX)
		asm.MULXQ(qInv0(b.elementName), regM, bavard.AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		asm.XORQ(bavard.AX, bavard.AX)

		// C,_ := t[0] + m*q[0]
		asm.MULXQ(qAt(0, b.elementName), bavard.AX, tmp)
		asm.ADCXQ(t[0], bavard.AX)
		asm.MOVQ(tmp, t[0])

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < b.nbWords; j++ {
			asm.ADCXQ(t[j], t[j-1])
			asm.MULXQ(qAt(j, b.elementName), bavard.AX, t[j])
			asm.ADOXQ(bavard.AX, t[j-1])
		}
		asm.MOVQ(0, bavard.AX)
		asm.ADCXQ(bavard.AX, t[b.nbWordsLastIndex])
		asm.ADOXQ(A, t[b.nbWordsLastIndex])

		asm.PushRegister(tmp)
	}

	// free registers
	asm.PushRegister(x, A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := asm.PopRegister()
	asm.MOVQ("res+0(FP)", r)
	b.reduce(asm, t, r)
	asm.RET()

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		asm.WriteLn("no_adx:")
		asm.Reset()
		x := asm.PopRegister()
		y := asm.PopRegister()
		// dereference x and y
		asm.MOVQ("x+8(FP)", x)
		asm.MOVQ("x+8(FP)", y)
		b.mulNoAdx(asm, x, y)
	}

	return nil
}
