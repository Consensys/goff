package cmd

import (
	"fmt"
)

func (b *asmBuilder) square(F *field) error {

	const header = `
// func square%s(res,y *%s)
TEXT ·square%s(SB), NOSPLIT, $0-16
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

	// if adx and mulx instructions are not available, uses MUL algorithm.
	`

	b.WriteLn(fmt.Sprintf(header, F.ElementName, F.ElementName, F.ElementName))

	// registers
	b.registers = make([]register, len(staticRegisters))
	copy(b.registers, staticRegisters) // re init registers in case

	regT := make([]register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
	}

	b.CMPB("·supportAdx(SB)", 1, "check if we support MULX and ADOX instructions")
	b.JNE("no_adx", "no support for MULX or ADOX instructions")

	{
		regY := b.PopRegister()
		b.MOVQ("y+8(FP)", regY, "dereference y")
		regA := b.PopRegister()

		for i := 0; i < F.NbWords; i++ {
			b.Comment(fmt.Sprintf("outter loop %d", i))
			b.XORQ(ax, ax, "clear up flags")

			b.Comment(fmt.Sprintf("dx = y[%d]", i))
			b.MOVQ(regY.at(i), dx)

			// instead of
			// for j=i+1 to N-1
			//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
			// we first add the x[j]*x[i] to a temporary u (set of registers)
			// set double it, before doing
			// for j=i+1 to N-1
			//     A,t[j] = u[j] + t[j] + A
			if i != F.NbWordsLastIndex {
				regU := make([]register, (F.NbWords - i - 1))
				for i := 0; i < len(regU); i++ {
					regU[i] = b.PopRegister()
				}
				offset := i + 1

				// 1- compute u = x[j] * x[i]
				// for j=i+1 to N-1
				//     A,u[j] = x[j]*x[i] + A
				if (i + 1) == F.NbWordsLastIndex {
					b.MULXQ(regY.at(i+1), regU[0], regA)
				} else {
					for j := i + 1; j < F.NbWords; j++ {
						if j == i+1 {
							// first iteration
							b.MULXQ(regY.at(j), regU[j-offset], regU[j+1-offset])
						} else {
							if j == F.NbWordsLastIndex {
								b.MULXQ(regY.at(j), ax, regA)
							} else {
								b.MULXQ(regY.at(j), ax, regU[j+1-offset])
							}
							b.ADCXQ(ax, regU[j-offset])
						}
					}
					b.MOVQ(0, ax)
					b.ADCXQ(ax, regA)
					b.XORQ(ax, ax, "clear up flags")
				}

				if i == 0 {
					// C, t[i] = x[i] * x[i] + t[i]
					b.MULXQ(dx, regT[i], dx)

					// when i == 0, T is not set yet
					// so  we can use ADOXQ carry chain to propagate C from x[i] * x[i] + t[i] (dx)

					// for j=i+1 to N-1
					// 		C, t[j] = u[j] + u[j] + t[j] + C
					for j := 0; j < len(regU); j++ {
						b.ADCXQ(regU[j], regU[j])
						b.MOVQ(regU[j], regT[j+offset])
						if j == 0 {
							b.ADOXQ(dx, regT[j+offset])
						} else {
							b.ADOXQ(ax, regT[j+offset])
						}
					}

					b.ADCXQ(regA, regA)
					b.ADOXQ(ax, regA)

				} else {
					// i != 0 so T is set.
					// we first use ADOXQ carry chain to perform t = u + u + t
					for j := 0; j < len(regU); j++ {
						b.ADCXQ(regU[j], regU[j])
						b.ADOXQ(regU[j], regT[j+offset])
					}

					b.ADCXQ(regA, regA)
					b.ADOXQ(ax, regA)

					// reset flags
					b.XORQ(ax, ax, "clear up flags")

					// C, t[i] = x[i] * x[i] + t[i]
					b.MULXQ(dx, ax, dx)
					b.ADOXQ(ax, regT[i])
					b.MOVQ(0, ax)

					// propagate C
					for j := i + 1; j < F.NbWords; j++ {
						if j == i+1 {
							b.ADOXQ(dx, regT[j])
						} else {
							b.ADOXQ(ax, regT[j])
						}
					}

					b.ADOXQ(ax, regA)
				}

				b.PushRegister(regU...)

			} else {
				// i == last index
				b.MULXQ(dx, ax, regA)
				b.ADCXQ(ax, regT[i])
				b.MOVQ(0, ax)
				b.ADCXQ(ax, regA)
			}

			regM := b.PopRegister()
			// m := t[0]*q'[0] mod W
			b.MOVQ(F.QInverse[0], dx)
			b.MULXQ(regT[0], regM, dx)

			// clear the carry flags
			b.XORQ(dx, dx, "clear up flags")

			// C,_ := t[0] + m*q[0]
			b.MOVQ(F.Q[0], dx)
			b.MULXQ(regM, ax, dx)
			b.ADCXQ(regT[0], ax)
			b.MOVQ(dx, regT[0])

			// for j=1 to N-1
			//    (C,t[j-1]) := t[j] + m*q[j] + C
			for j := 1; j < F.NbWords; j++ {
				b.MOVQ(F.Q[j], dx)
				b.ADCXQ(regT[j], regT[j-1])
				b.MULXQ(regM, ax, regT[j])
				b.ADOXQ(ax, regT[j-1])
			}
			b.MOVQ(0, ax)
			b.ADCXQ(ax, regT[F.NbWordsLastIndex])
			b.ADOXQ(regA, regT[F.NbWordsLastIndex])

			b.PushRegister(regM)
		}

		// free registers

		b.PushRegister(regY, regA)
	}

	// ---------------------------------------------------------------------------------------------
	// reduce
	regX := b.PopRegister()
	b.Comment("dereference res")
	b.MOVQ("res+0(FP)", regX)
	b.reduce(F, regT, regX)

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		b.WriteLn("no_adx:")

		regC := b.PopRegister()
		regYi := b.PopRegister()
		regA := b.PopRegister()
		regM := b.PopRegister()
		regY := b.PopRegister()

		b.Comment("dereference y")
		b.MOVQ("y+8(FP)", regY)

		for i := 0; i < F.NbWords; i++ {
			// (A,t[0]) := t[0] + x[0]*y[{{$i}}]
			b.MOVQ(regY.at(0), ax)
			b.MOVQ(regY.at(i), regYi)
			b.MULQ(regYi)
			if i != 0 {
				b.ADDQ(ax, regT[0])
				b.ADCQ(0, dx)
			} else {
				b.MOVQ(ax, regT[0])
			}
			b.MOVQ(dx, regA)

			// m := t[0]*q'[0] mod W
			b.MOVQ(F.QInverse[0], regM)
			b.IMULQ(regT[0], regM)

			// C,_ := t[0] + m*q[0]
			b.MOVQ(F.Q[0], ax)
			b.MULQ(regM)
			b.ADDQ(regT[0], ax)
			b.ADCQ(0, dx)
			b.MOVQ(dx, regC)

			// for j=1 to N-1
			//    (A,t[j])  := t[j] + x[j]*y[i] + A
			//    (C,t[j-1]) := t[j] + m*q[j] + C
			for j := 1; j < F.NbWords; j++ {
				b.MOVQ(regY.at(j), ax)
				b.MULQ(regYi)
				if i != 0 {
					b.ADDQ(regA, regT[j])
					b.ADCQ(0, dx)
					b.ADDQ(ax, regT[j])
					b.ADCQ(0, dx)
				} else {
					b.MOVQ(regA, regT[j])
					b.ADDQ(ax, regT[j])
					b.ADCQ(0, dx)
				}
				b.MOVQ(dx, regA)

				b.MOVQ(F.Q[j], ax)
				b.MULQ(regM)
				b.ADDQ(regT[j], regC)
				b.ADCQ(0, dx)
				b.ADDQ(ax, regC)
				b.ADCQ(0, dx)
				b.MOVQ(regC, regT[j-1])
				b.MOVQ(dx, regC)
			}

			b.ADDQ(regC, regA)
			b.MOVQ(regA, regT[F.NbWordsLastIndex])

		}
		b.Comment("dereference res")
		b.MOVQ("res+0(FP)", regX)
		b.JMP("reduce")
	}

	return nil
}
