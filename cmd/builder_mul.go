package cmd

import (
	"fmt"
)

type mulType uint8

const (
	mulAssign mulType = iota
	fromMont
)

func (b *asmBuilder) mulNoCarry(F *field, mType mulType) error {

	var header string
	switch mType {
	case mulAssign:
		header = mulHeader
	case fromMont:
		header = fromMontHeader
	}

	b.WriteLn(fmt.Sprintf(header, F.ElementName, F.ElementName, F.ElementName))

	// registers
	b.registers = make([]register, len(staticRegisters))
	copy(b.registers, staticRegisters) // re init registers in case
	var regX register

	regT := make([]register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
	}

	regX = b.PopRegister()
	b.Comment("dereference x")
	b.MOVQ("res+0(FP)", regX)

	if mType == fromMont {
		// 	for i=0 to N-1
		//     t[i] = a[i]
		for i := 0; i < F.NbWords; i++ {
			b.MOVQ(regX.at(i), regT[i])
		}
	}

	b.Comment("check if we support adx and mulx")
	b.CMPB("·supportAdx(SB)", 1)
	b.JNE("no_adx")

	{

		regM := b.PopRegister()
		var regA, regY register
		if mType != fromMont {
			regA = b.PopRegister()
			regY = b.PopRegister()
			b.Comment("dereference y")
			b.MOVQ("y+8(FP)", regY)
		}

		for i := 0; i < F.NbWords; i++ {
			b.Comment(fmt.Sprintf("outter loop %d", i))
			b.Comment("clear up flags")
			b.XORQ(dx, dx)

			if mType != fromMont {
				b.Comment(fmt.Sprintf("DX = y[%d]", i))
				b.MOVQ(regY.at(i), dx)

				// for j=0 to N-1
				//    (A,t[j])  := t[j] + x[j]*y[i] + A
				for j := 0; j < F.NbWords; j++ {
					reg := regA
					if i == 0 {
						if j == 0 {
							b.MULXQ(regX.at(j), regT[j], regT[j+1])
						} else if j != F.NbWordsLastIndex {
							reg = regT[j+1]
						}
					} else if j != 0 {
						b.ADCXQ(regA, regT[j])
					}

					if !(i == 0 && j == 0) {
						b.MULXQ(regX.at(j), ax, reg)
						b.ADOXQ(ax, regT[j])
					}
				}

				// add the last carries to {{reg $regA}}
				b.MOVQ(0, dx)
				b.ADCXQ(dx, regA)
				b.ADOXQ(dx, regA)
			}

			// m := t[0]*q'[0] mod W
			b.MOVQ(F.QInverse[0], dx)
			b.MULXQ(regT[0], regM, dx)

			// clear the carry flags
			b.XORQ(dx, dx)

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
			if mType != fromMont {
				b.ADOXQ(regA, regT[F.NbWordsLastIndex])
			} else {
				b.ADOXQ(ax, regT[F.NbWordsLastIndex])
			}

		}

		// free registers
		b.PushRegister(regM)
		if mType != fromMont {
			b.PushRegister(regY, regA)
		}
	}

	// ---------------------------------------------------------------------------------------------
	// reduce
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

		if mType != fromMont {
			b.Comment("dereference y")
			b.MOVQ("y+8(FP)", regY)
		}

		for i := 0; i < F.NbWords; i++ {
			if mType != fromMont {
				// (A,t[0]) := t[0] + x[0]*y[{{$i}}]
				b.MOVQ(regX.at(0), ax)
				b.MOVQ(regY.at(i), regYi)
				b.MULQ(regYi)
				if i != 0 {
					b.ADDQ(ax, regT[0])
					b.ADCQ(0, dx)
				} else {
					b.MOVQ(ax, regT[0])
				}
				b.MOVQ(dx, regA)
			}

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
				if mType != fromMont {
					b.MOVQ(regX.at(j), ax)
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
				}

				b.MOVQ(F.Q[j], ax)
				b.MULQ(regM)
				b.ADDQ(regT[j], regC)
				b.ADCQ(0, dx)
				b.ADDQ(ax, regC)
				b.ADCQ(0, dx)
				b.MOVQ(regC, regT[j-1])
				b.MOVQ(dx, regC)
			}

			if mType != fromMont {
				b.ADDQ(regC, regA)
				b.MOVQ(regA, regT[F.NbWordsLastIndex])
			} else {
				b.MOVQ(regC, regT[F.NbWordsLastIndex])
			}

		}

		b.JMP("reduce")
	}

	return nil
}

func (b *asmBuilder) reduce(F *field, regT []register, result register) error {
	b.WriteLn("reduce:")

	// u = t - q
	regU := make([]register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regU[i] = b.PopRegister()
		b.MOVQ(regT[i], regU[i])
		b.MOVQ(F.Q[i], dx)

		if i == 0 {
			b.SUBQ(dx, regU[i])
		} else {
			b.SBBQ(dx, regU[i])
		}
	}
	// no borrow we return t
	b.JCS("t_is_smaller")

	// return u
	for i := 0; i < F.NbWords; i++ {
		b.MOVQ(regU[i], result.at(i))
	}
	b.RET()

	// return t
	b.WriteLn("t_is_smaller:")
	for i := 0; i < F.NbWords; i++ {
		b.MOVQ(regT[i], result.at(i))
	}
	b.RET()

	b.PushRegister(regU...)
	return nil
}

const mulHeader = `

// func mulAssign%s(res,y *%s)
// montgomery multiplication of res by y 
// stores the result in res
TEXT ·mulAssign%s(SB), NOSPLIT, $0-16
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
`

const fromMontHeader = `

// func fromMont%s(res *%s)
// montgomery multiplication of res by 1 
// stores the result in res
TEXT ·fromMont%s(SB), NOSPLIT, $0-8
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
	// 		t[N-1] = C

`
