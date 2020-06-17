package asm

import (
	"github.com/consensys/bavard"
)

func (b *Builder) sub(asm *bavard.Assembly) error {
	asm.FuncHeader("sub"+b.elementName, 0, 24)

	// registers
	t := asm.PopRegisters(b.nbWords)
	r := asm.PopRegister()

	// set DX to 0
	asm.XORQ(bavard.DX, bavard.DX)

	// dereference x
	asm.MOVQ("x+8(FP)", r)

	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(r.At(i), t[i])
	}

	// dereference y
	asm.MOVQ("y+16(FP)", r)

	// z = x - y mod q
	// move t = x
	asm.SUBQ(r.At(0), t[0])
	for i := 1; i < b.nbWords; i++ {
		asm.SBBQ(r.At(i), t[i])
	}

	// move modulus q into set of registers (overwrite with 0 if borrow is set)
	q := asm.PopRegisters(b.nbWords)
	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(b.q[i], q[i])
		asm.CMOVQCC(bavard.DX, q[i])
	}

	// dereference result
	asm.MOVQ("res+0(FP)", r)

	// add registers (q or 0) to t, and set to result
	asm.ADDQ(q[0], t[0])
	asm.MOVQ(t[0], r.At(0))
	for i := 1; i < b.nbWords; i++ {
		asm.ADCQ(q[i], t[i])
		asm.MOVQ(t[i], r.At(i))
	}

	asm.RET()

	return nil
}
