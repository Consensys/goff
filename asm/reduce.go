package asm

import (
	"fmt"

	"github.com/consensys/bavard"
)

func (b *builder) reduceFn(asm *bavard.Assembly) error {
	stackSize := 0
	if b.nbWords > smallModulus {
		stackSize = b.nbWords * 8
	}
	asm.FuncHeader("reduce"+b.elementName, stackSize, 8)

	r := asm.PopRegister()
	asm.MOVQ("res+0(FP)", r)

	t := asm.PopRegisters(b.nbWords)
	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(r.At(i), t[i])
	}

	b.reduce(asm, t, r)

	asm.RET()

	return nil
}

func (b *builder) reduce(asm *bavard.Assembly, t []bavard.Register, result bavard.Register) error {
	if b.nbWords > smallModulus {
		return b.reduceLarge(asm, t, result)
	}
	// u = t - q
	regU := asm.PopRegisters(b.nbWords)

	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(t[i], regU[i])

		if i == 0 {
			asm.SUBQ(qAt(i, b.elementName), regU[i])
		} else {
			asm.SBBQ(qAt(i, b.elementName), regU[i])
		}
	}

	// conditional move of u into t (if we have a borrow we need to return t - q)
	for i := 0; i < b.nbWords; i++ {
		asm.CMOVQCC(regU[i], t[i])
	}

	// return t
	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(t[i], result.At(i))
	}

	asm.PushRegister(regU...)
	return nil
}

func (b *builder) reduceLarge(asm *bavard.Assembly, t []bavard.Register, result bavard.Register) error {
	// u = t - q
	u := make([]string, b.nbWords)

	for i := 0; i < b.nbWords; i++ {
		// use stack
		u[i] = fmt.Sprintf("t%d-%d(SP)", i, 8+i*8)
		asm.MOVQ(t[i], u[i])

		if i == 0 {
			asm.SUBQ(qAt(i, b.elementName), t[i])
		} else {
			asm.SBBQ(qAt(i, b.elementName), t[i])
		}
	}

	// conditional move of u into t (if we have a borrow we need to return t - q)
	for i := 0; i < b.nbWords; i++ {
		asm.CMOVQCS(u[i], t[i])
	}

	// return t
	for i := 0; i < b.nbWords; i++ {
		asm.MOVQ(t[i], result.At(i))
	}

	return nil
}
