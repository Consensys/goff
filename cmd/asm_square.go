package cmd

import (
	"fmt"

	"github.com/consensys/bavard"
)

func generateSquareASM(b *bavard.Assembly, F *field) error {

	b.FuncHeader("square"+F.ElementName, 16)

	b.WriteLn(`
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
	`)

	// registers
	b.Reset()

	regT := make([]bavard.Register, F.NbWords)
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
			b.XORQ(bavard.AX, bavard.AX, "clear up flags")

			b.Comment(fmt.Sprintf("dx = y[%d]", i))
			b.MOVQ(regY.At(i), bavard.DX)

			// instead of
			// for j=i+1 to N-1
			//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
			// we first add the x[j]*x[i] to a temporary u (set of registers)
			// set double it, before doing
			// for j=i+1 to N-1
			//     A,t[j] = u[j] + t[j] + A
			if i != F.NbWordsLastIndex {
				regU := make([]bavard.Register, (F.NbWords - i - 1))
				for i := 0; i < len(regU); i++ {
					regU[i] = b.PopRegister()
				}
				offset := i + 1

				// 1- compute u = x[j] * x[i]
				// for j=i+1 to N-1
				//     A,u[j] = x[j]*x[i] + A
				if (i + 1) == F.NbWordsLastIndex {
					b.MULXQ(regY.At(i+1), regU[0], regA)
				} else {
					for j := i + 1; j < F.NbWords; j++ {
						yj := regY.At(j)
						if j == i+1 {
							// first iteration
							b.MULXQ(yj, regU[j-offset], regU[j+1-offset])
						} else {
							if j == F.NbWordsLastIndex {
								b.MULXQ(yj, bavard.AX, regA)
							} else {
								b.MULXQ(yj, bavard.AX, regU[j+1-offset])
							}
							b.ADCXQ(bavard.AX, regU[j-offset])
						}
					}
					b.MOVQ(0, bavard.AX)
					b.ADCXQ(bavard.AX, regA)
					b.XORQ(bavard.AX, bavard.AX, "clear up flags")
				}

				if i == 0 {
					// C, t[i] = x[i] * x[i] + t[i]
					b.MULXQ(bavard.DX, regT[i], bavard.DX)

					// when i == 0, T is not set yet
					// so  we can use ADOXQ carry chain to propagate C from x[i] * x[i] + t[i] (dx)

					// for j=i+1 to N-1
					// 		C, t[j] = u[j] + u[j] + t[j] + C
					for j := 0; j < len(regU); j++ {
						b.ADCXQ(regU[j], regU[j])
						b.MOVQ(regU[j], regT[j+offset])
						if j == 0 {
							b.ADOXQ(bavard.DX, regT[j+offset])
						} else {
							b.ADOXQ(bavard.AX, regT[j+offset])
						}
					}

					b.ADCXQ(regA, regA)
					b.ADOXQ(bavard.AX, regA)

				} else {
					// i != 0 so T is set.
					// we first use ADOXQ carry chain to perform t = u + u + t
					for j := 0; j < len(regU); j++ {
						b.ADCXQ(regU[j], regU[j])
						b.ADOXQ(regU[j], regT[j+offset])
					}

					b.ADCXQ(regA, regA)
					b.ADOXQ(bavard.AX, regA)

					// reset flags
					b.XORQ(bavard.AX, bavard.AX, "clear up flags")

					// C, t[i] = x[i] * x[i] + t[i]
					b.MULXQ(bavard.DX, bavard.AX, bavard.DX)
					b.ADOXQ(bavard.AX, regT[i])
					b.MOVQ(0, bavard.AX)

					// propagate C
					for j := i + 1; j < F.NbWords; j++ {
						if j == i+1 {
							b.ADOXQ(bavard.DX, regT[j])
						} else {
							b.ADOXQ(bavard.AX, regT[j])
						}
					}

					b.ADOXQ(bavard.AX, regA)
				}

				b.PushRegister(regU...)

			} else {
				// i == last index
				b.MULXQ(bavard.DX, bavard.AX, regA)
				b.ADCXQ(bavard.AX, regT[i])
				b.MOVQ(0, bavard.AX)
				b.ADCXQ(bavard.AX, regA)
			}

			regTmp := b.PopRegister()
			// m := t[0]*q'[0] mod W
			regM := bavard.DX
			b.MOVQ(regT[0], bavard.DX)
			b.MULXQ(qInv0(F), regM, bavard.AX, "m := t[0]*q'[0] mod W")
			// b.MOVQ(F.QInverse[0], bavard.DX)
			// b.MULXQ(regT[0], regM, bavard.DX)

			// clear the carry flags
			b.XORQ(bavard.AX, bavard.AX, "clear up flags")

			// C,_ := t[0] + m*q[0]
			b.MULXQ(qAt(0, F), bavard.AX, regTmp)
			// b.MOVQ(F.Q[0], bavard.DX)
			// b.MULXQ(regM, bavard.AX, bavard.DX)
			b.ADCXQ(regT[0], bavard.AX)
			b.MOVQ(regTmp, regT[0])

			// for j=1 to N-1
			//    (C,t[j-1]) := t[j] + m*q[j] + C
			for j := 1; j < F.NbWords; j++ {
				// b.MOVQ(F.Q[j], bavard.DX)
				b.ADCXQ(regT[j], regT[j-1])
				b.MULXQ(qAt(j, F), bavard.AX, regT[j])
				b.ADOXQ(bavard.AX, regT[j-1])
			}
			b.MOVQ(0, bavard.AX)
			b.ADCXQ(bavard.AX, regT[F.NbWordsLastIndex])
			b.ADOXQ(regA, regT[F.NbWordsLastIndex])

			b.PushRegister(regTmp)
		}

		// free registers

		b.PushRegister(regY, regA)
	}

	// ---------------------------------------------------------------------------------------------
	// reduce
	regX := b.PopRegister()
	b.Comment("dereference res")
	b.MOVQ("res+0(FP)", regX)
	generateReduceASM(b, F, regT, regX)

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
			b.MOVQ(regY.At(0), bavard.AX)
			b.MOVQ(regY.At(i), regYi)
			b.MULQ(regYi)
			if i != 0 {
				b.ADDQ(bavard.AX, regT[0])
				b.ADCQ(0, bavard.DX)
			} else {
				b.MOVQ(bavard.AX, regT[0])
			}
			b.MOVQ(bavard.DX, regA)

			// m := t[0]*q'[0] mod W
			b.MOVQ(F.QInverse[0], regM)
			b.IMULQ(regT[0], regM)

			// C,_ := t[0] + m*q[0]
			b.MOVQ(F.Q[0], bavard.AX)
			b.MULQ(regM)
			b.ADDQ(regT[0], bavard.AX)
			b.ADCQ(0, bavard.DX)
			b.MOVQ(bavard.DX, regC)

			// for j=1 to N-1
			//    (A,t[j])  := t[j] + x[j]*y[i] + A
			//    (C,t[j-1]) := t[j] + m*q[j] + C
			for j := 1; j < F.NbWords; j++ {
				b.MOVQ(regY.At(j), bavard.AX)
				b.MULQ(regYi)
				if i != 0 {
					b.ADDQ(regA, regT[j])
					b.ADCQ(0, bavard.DX)
					b.ADDQ(bavard.AX, regT[j])
					b.ADCQ(0, bavard.DX)
				} else {
					b.MOVQ(regA, regT[j])
					b.ADDQ(bavard.AX, regT[j])
					b.ADCQ(0, bavard.DX)
				}
				b.MOVQ(bavard.DX, regA)

				b.MOVQ(F.Q[j], bavard.AX)
				b.MULQ(regM)
				b.ADDQ(regT[j], regC)
				b.ADCQ(0, bavard.DX)
				b.ADDQ(bavard.AX, regC)
				b.ADCQ(0, bavard.DX)
				b.MOVQ(regC, regT[j-1])
				b.MOVQ(bavard.DX, regC)
			}

			b.ADDQ(regC, regA)
			b.MOVQ(regA, regT[F.NbWordsLastIndex])

		}
		b.Comment("dereference res")
		b.PushRegister(regC, regYi, regA, regM, regY)
		b.MOVQ("res+0(FP)", regX)
		generateReduceASM(b, F, regT, regX)
	}

	return nil
}
