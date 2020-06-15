package cmd

import (
	"fmt"

	"github.com/consensys/bavard"
)

func qAt(index int, F *field) string {
	return fmt.Sprintf("·q%s+%d(SB)", F.ElementName, index*8)
}

func qInv0(F *field) string {
	return fmt.Sprintf("·q%sInv0(SB)", F.ElementName)
}

func generateFromMontASM(b *bavard.Assembly, F *field) error {
	// reset register state
	b.Reset()

	b.FuncHeader("fromMont"+F.ElementName, 8)
	b.WriteLn(`
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

	regT := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
	}

	regX := b.PopRegister()
	b.MOVQ("res+0(FP)", regX, "dereference x")

	// 	for i=0 to N-1
	//     t[i] = a[i]
	for i := 0; i < F.NbWords; i++ {
		b.MOVQ(regX.At(i), regT[i], fmt.Sprintf("t[%d] = x[%d]", i, i))
	}

	b.CMPB("·supportAdx(SB)", 1, "check if we support MULX and ADOX instructions")
	b.JNE("no_adx", "no support for MULX or ADOX instructions")

	{
		regTmp := b.PopRegister()

		for i := 0; i < F.NbWords; i++ {
			b.Comment(fmt.Sprintf("outter loop %d", i))
			b.XORQ(bavard.DX, bavard.DX, "clear up flags")

			// m := t[0]*q'[0] mod W
			regM := bavard.DX
			b.MOVQ(regT[0], bavard.DX)
			b.MULXQ(qInv0(F), regM, bavard.AX, "m := t[0]*q'[0] mod W")

			// clear the carry flags
			b.XORQ(bavard.AX, bavard.AX, "clear the flags")

			// C,_ := t[0] + m*q[0]
			b.Comment("C,_ := t[0] + m*q[0]")

			b.MULXQ(qAt(0, F), bavard.AX, regTmp)
			b.ADCXQ(regT[0], bavard.AX)
			b.MOVQ(regTmp, regT[0])

			b.Comment("for j=1 to N-1")
			b.Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

			// for j=1 to N-1
			//    (C,t[j-1]) := t[j] + m*q[j] + C
			for j := 1; j < F.NbWords; j++ {
				b.ADCXQ(regT[j], regT[j-1])
				b.MULXQ(qAt(j, F), bavard.AX, regT[j])
				b.ADOXQ(bavard.AX, regT[j-1])
			}
			b.MOVQ(0, bavard.AX)
			b.ADCXQ(bavard.AX, regT[F.NbWordsLastIndex])
			b.ADOXQ(bavard.AX, regT[F.NbWordsLastIndex])

		}

		// free registers
		b.PushRegister(regTmp)
	}

	// ---------------------------------------------------------------------------------------------
	// reduce
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

		for i := 0; i < F.NbWords; i++ {
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
				b.MOVQ(F.Q[j], bavard.AX)
				b.MULQ(regM)
				b.ADDQ(regT[j], regC)
				b.ADCQ(0, bavard.DX)
				b.ADDQ(bavard.AX, regC)
				b.ADCQ(0, bavard.DX)
				b.MOVQ(regC, regT[j-1])
				b.MOVQ(bavard.DX, regC)
			}

			b.MOVQ(regC, regT[F.NbWordsLastIndex])

		}

		// ---------------------------------------------------------------------------------------------
		// reduce
		b.PushRegister(regC, regYi, regA, regM, regY)
		generateReduceASM(b, F, regT, regX)
	}

	return nil
}

func generateMulASM(b *bavard.Assembly, F *field) error {
	// reset register state
	b.Reset()

	b.FuncHeader("mul"+F.ElementName, 24)

	b.WriteLn(`
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

	regT := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
	}

	regX := b.PopRegister()
	b.MOVQ("x+8(FP)", regX, "dereference x")

	b.CMPB("·supportAdx(SB)", 1, "check if we support MULX and ADOX instructions")
	b.JNE("no_adx", "no support for MULX or ADOX instructions")

	{
		regTmp := b.PopRegister()
		var regxi []bavard.Register
		regA := b.PopRegister()
		regY := b.PopRegister()
		b.MOVQ("y+16(FP)", regY, "dereference y")

		cacheSize := b.AvailableRegisters()
		if cacheSize > F.NbWords {
			cacheSize = F.NbWords
		}
		regxi = make([]bavard.Register, cacheSize)
		for i := 0; i < len(regxi); i++ {
			regxi[i] = b.PopRegister()
			b.MOVQ(regX.At(i), regxi[i], fmt.Sprintf("%s = x[%d]", string(regxi[i]), i))
		}

		for i := 0; i < F.NbWords; i++ {
			b.Comment(fmt.Sprintf("outter loop %d", i))
			b.XORQ(bavard.DX, bavard.DX, "clear up flags")

			b.MOVQ(regY.At(i), bavard.DX, fmt.Sprintf("DX = y[%d]", i))
			// for j=0 to N-1
			//    (A,t[j])  := t[j] + x[j]*y[i] + A
			for j := 0; j < F.NbWords; j++ {
				xj := regX.At(j)
				if j < len(regxi) {
					xj = string(regxi[j])
				}

				reg := regA
				if i == 0 {
					if j == 0 {
						b.MULXQ(xj, regT[j], regT[j+1], fmt.Sprintf("t[%d], t[%d] = y[%d] * x[%d]", j, j+1, i, j))
					} else if j != F.NbWordsLastIndex {
						reg = regT[j+1]
					}
				} else if j != 0 {
					b.ADCXQ(regA, regT[j], fmt.Sprintf("t[%d] += regA", j))
				}

				if !(i == 0 && j == 0) {
					b.MULXQ(xj, bavard.AX, reg)
					b.ADOXQ(bavard.AX, regT[j])
				}
			}

			b.Comment("add the last carries to " + string(regA))
			b.MOVQ(0, bavard.DX)
			b.ADCXQ(bavard.DX, regA)
			b.ADOXQ(bavard.DX, regA)

			// m := t[0]*q'[0] mod W
			regM := bavard.DX
			b.MOVQ(regT[0], bavard.DX)
			b.MULXQ(qInv0(F), regM, bavard.AX, "m := t[0]*q'[0] mod W")

			// clear the carry flags
			b.XORQ(bavard.AX, bavard.AX, "clear the flags")

			// C,_ := t[0] + m*q[0]
			b.Comment("C,_ := t[0] + m*q[0]")

			b.MULXQ(qAt(0, F), bavard.AX, regTmp)
			b.ADCXQ(regT[0], bavard.AX)
			b.MOVQ(regTmp, regT[0])

			b.Comment("for j=1 to N-1")
			b.Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

			// for j=1 to N-1
			//    (C,t[j-1]) := t[j] + m*q[j] + C
			for j := 1; j < F.NbWords; j++ {
				b.ADCXQ(regT[j], regT[j-1])
				b.MULXQ(qAt(j, F), bavard.AX, regT[j])
				b.ADOXQ(bavard.AX, regT[j-1])
			}
			b.MOVQ(0, bavard.AX)
			b.ADCXQ(bavard.AX, regT[F.NbWordsLastIndex])
			b.ADOXQ(regA, regT[F.NbWordsLastIndex])
		}

		// free registers
		b.PushRegister(regTmp)
		b.PushRegister(regY, regA)
		b.PushRegister(regxi...)
	}

	// ---------------------------------------------------------------------------------------------
	// reduce
	b.MOVQ("res+0(FP)", regX, "dereference res")
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

		b.MOVQ("y+16(FP)", regY, "dereference y")

		for i := 0; i < F.NbWords; i++ {
			// (A,t[0]) := t[0] + x[0]*y[{{$i}}]
			b.MOVQ(regX.At(0), bavard.AX)
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
				b.MOVQ(regX.At(j), bavard.AX)
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

		// ---------------------------------------------------------------------------------------------
		// reduce
		b.PushRegister(regC, regYi, regA, regM, regY)
		b.MOVQ("res+0(FP)", regX, "dereference res")
		generateReduceASM(b, F, regT, regX)
	}

	return nil
}

func generateReduceFuncASM(b *bavard.Assembly, F *field) error {
	// reset register state
	b.Reset()

	b.FuncHeader("reduce"+F.ElementName, 8)

	regX := b.PopRegister()
	b.MOVQ("res+0(FP)", regX, "dereference x")

	regT := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
		b.MOVQ(regX.At(i), regT[i], fmt.Sprintf("t[%d] = x[%d]", i, i))
	}

	return generateReduceASM(b, F, regT, regX)
}

func generateReduceASM(b *bavard.Assembly, F *field, regT []bavard.Register, result bavard.Register) error {
	// u = t - q
	regU := make([]bavard.Register, F.NbWords)

	for i := 0; i < F.NbWords; i++ {
		regU[i] = b.PopRegister()
		b.MOVQ(regT[i], regU[i])

		if i == 0 {
			b.SUBQ(qAt(i, F), regU[i])
		} else {
			b.SBBQ(qAt(i, F), regU[i])
		}
	}

	// conditional move of u into t (if we have a borrow we need to return t - q)
	for i := 0; i < F.NbWords; i++ {
		b.CMOVQCC(regU[i], regT[i])
	}

	// return t
	for i := 0; i < F.NbWords; i++ {
		b.MOVQ(regT[i], result.At(i))
	}
	b.RET()

	b.PushRegister(regU...)
	return nil
}
