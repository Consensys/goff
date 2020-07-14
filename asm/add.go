package asm

import (
	"github.com/consensys/bavard"
)

func (b *builder) double(asm *bavard.Assembly) error {
	// func header
	stackSize := 0
	if b.nbWords > smallModulus {
		stackSize = b.nbWords * 8
	}
	asm.FuncHeader("double"+b.elementName, stackSize, 16)

	// registers
	r := bavard.Register(bavard.AX)

	// dereference x
	asm.MOVQ("x+8(FP)", r)

	// move t = x
	t := asm.PopRegisters(b.nbWords)
	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(r.At(i), t[i])
	}

	// t = t + y = x + y
	asm.ADDQ(t[0], t[0])
	for i := 1; i < b.nbWords; i++ {
		asm.ADCQ(t[i], t[i])
	}

	asm.Comment("note that we don't check for the carry here, as this code was generated assuming F.NoCarry condition is set")
	asm.Comment("(see goff for more details)")

	// reduce
	asm.MOVQ("res+0(FP)", r)

	b.reduce(asm, t, r)

	asm.RET()

	return nil
}

func (b *builder) add(asm *bavard.Assembly) error {
	// func header
	stackSize := 0
	if b.nbWords > smallModulus {
		stackSize = b.nbWords * 8
	}
	asm.FuncHeader("add"+b.elementName, stackSize, 24)

	// registers
	r := bavard.Register(bavard.AX)

	// dereference x
	asm.MOVQ("x+8(FP)", r)

	// move t = x
	t := asm.PopRegisters(b.nbWords)
	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(r.At(i), t[i])
	}

	// dereference y
	asm.MOVQ("y+16(FP)", r)

	// t = t + y = x + y
	asm.ADDQ(r.At(0), t[0])
	for i := 1; i < b.nbWords; i++ {
		asm.ADCQ(r.At(i), t[i])
	}

	asm.Comment("note that we don't check for the carry here, as this code was generated assuming F.NoCarry condition is set")
	asm.Comment("(see goff for more details)")

	// reduce
	asm.MOVQ("res+0(FP)", r)

	b.reduce(asm, t, r)

	asm.RET()

	return nil
}
