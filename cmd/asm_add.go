package cmd

import (
	"fmt"

	"github.com/consensys/bavard"
)

type addType uint8

const (
	add addType = iota
	addAssign
	double
)

func generateAddASM(b *bavard.Assembly, F *field, aType addType) error {

	switch aType {
	case add:
		b.FuncHeader("add"+F.ElementName, 24)
	case addAssign:
		b.FuncHeader("addAssign"+F.ElementName, 16)
	case double:
		b.FuncHeader("double"+F.ElementName, 16)
	}

	// registers
	b.Reset()
	var regX bavard.Register

	regT := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
	}

	regX = b.PopRegister()
	regY := b.PopRegister()
	switch aType {
	case add:
		b.MOVQ("x+8(FP)", regX, "dereference x")
		b.MOVQ("y+16(FP)", regY, "dereference y")
	case addAssign, double:
		b.MOVQ("res+0(FP)", regX, "dereference x")
		b.MOVQ("y+8(FP)", regY, "dereference y")
	}

	if aType == double {
		for i := 0; i < F.NbWords; i++ {
			b.MOVQ(regY.At(i), regT[i], fmt.Sprintf("t[%d] = y[%d]", i, i))
		}
		b.ADDQ(regT[0], regT[0])
		for i := 1; i < F.NbWords; i++ {
			b.ADCQ(regT[i], regT[i])
		}
	} else {
		for i := 0; i < F.NbWords; i++ {
			b.MOVQ(regX.At(i), regT[i], fmt.Sprintf("t[%d] = x[%d]", i, i))
		}

		b.ADDQ(regY.At(0), regT[0])
		for i := 1; i < F.NbWords; i++ {
			b.ADCQ(regY.At(i), regT[i])
		}
	}

	// note: here we don't check the carry because we assume F.NoCarry is met to generated ASM code.
	// if moduli uses all bits, then we need a different code path.
	b.PushRegister(regY)
	if aType == add {
		b.MOVQ("res+0(FP)", regX, "dereference res")
	}

	generateReduceASM(b, F, regT, regX)

	return nil
}
