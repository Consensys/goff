package asm

import (
	"github.com/consensys/bavard"
)

func (b *Builder) fromMont(asm *bavard.Assembly) error {
	stackSize := 0
	if b.nbWords > SmallModulus {
		stackSize = b.nbWords * 8
	}
	asm.FuncHeader("_fromMontADX"+b.elementName, stackSize, 8)
	asm.WriteLn(`
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

	// registers
	t := asm.PopRegisters(b.nbWords)
	r := asm.PopRegister()

	// dereference y
	asm.MOVQ("res+0(FP)", r)

	// 	for i=0 to N-1
	//     t[i] = a[i]
	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(r.At(i), t[i])
	}

	var tmp bavard.Register
	if b.nbWords > 11 {
		tmp = r
	} else {
		tmp = asm.PopRegister()
	}
	for i := 0; i < b.nbWords; i++ {

		asm.XORQ(bavard.DX, bavard.DX)

		// m := t[0]*q'[0] mod W
		regM := bavard.DX
		asm.MOVQ(t[0], bavard.DX)
		asm.MULXQ(qInv0(b.elementName), regM, bavard.AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		asm.XORQ(bavard.AX, bavard.AX)

		// C,_ := t[0] + m*q[0]
		asm.Comment("C,_ := t[0] + m*q[0]")

		asm.MULXQ(qAt(0, b.elementName), bavard.AX, tmp)
		asm.ADCXQ(t[0], bavard.AX)
		asm.MOVQ(tmp, t[0])

		asm.Comment("for j=1 to N-1")
		asm.Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < b.nbWords; j++ {
			asm.ADCXQ(t[j], t[j-1])
			asm.MULXQ(qAt(j, b.elementName), bavard.AX, t[j])
			asm.ADOXQ(bavard.AX, t[j-1])
		}
		asm.MOVQ(0, bavard.AX)
		asm.ADCXQ(bavard.AX, t[b.nbWordsLastIndex])
		asm.ADOXQ(bavard.AX, t[b.nbWordsLastIndex])

	}

	if b.nbWords > 11 {
		asm.MOVQ("res+0(FP)", r)
	} else {
		asm.PushRegister(tmp)
	}
	// ---------------------------------------------------------------------------------------------
	// reduce
	b.reduce(asm, t, r)
	asm.RET()
	return nil
}
