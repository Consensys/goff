package asm

import (
	"fmt"

	"github.com/consensys/bavard"
)

func (b *builder) fromMont(asm *bavard.Assembly) error {
	stackSize := 8
	if b.nbWords > smallModulus {
		stackSize = b.nbWords * 8
	}
	asm.FuncHeader("_fromMontADX"+b.elementName, stackSize, 8)
	asm.WriteLn("NO_LOCAL_POINTERS")
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

	// check ADX instruction support
	asm.CMPB("·supportAdx(SB)", 1)
	asm.JNE("no_adx")

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
	hasRegisters := asm.AvailableRegisters() > 0
	if !hasRegisters {
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

	if !hasRegisters {
		asm.MOVQ("res+0(FP)", r)
	} else {
		asm.PushRegister(tmp)
	}
	// ---------------------------------------------------------------------------------------------
	// reduce
	b.reduce(asm, t, r)
	asm.RET()

	// No adx
	asm.WriteLn("no_adx:")
	asm.MOVQ("res+0(FP)", bavard.AX)
	asm.MOVQ(bavard.AX, "(SP)")
	asm.WriteLn(fmt.Sprintf("CALL ·_fromMontGeneric%s(SB)", b.elementName))
	asm.RET()
	return nil
}
