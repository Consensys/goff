package cmd

import (
	"fmt"

	"github.com/consensys/bavard"
)

func generateAddASM(b *bavard.Assembly, F *field) error {
	// reset register state
	b.Reset()

	b.FuncHeader("add"+F.ElementName, 24)

	regX := b.PopRegister()
	regY := b.PopRegister()
	b.MOVQ("x+8(FP)", regX, "dereference x")
	b.MOVQ("y+16(FP)", regY, "dereference y")

	// move x into T[]
	regT := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
		b.MOVQ(regX.At(i), regT[i], fmt.Sprintf("t[%d] = x[%d]", i, i))
	}

	// t = t + y = x + y
	b.ADDQ(regY.At(0), regT[0])
	for i := 1; i < F.NbWords; i++ {
		b.ADCQ(regY.At(i), regT[i])
	}
	b.Comment("note that we don't check for the carry here, as this code was generated assuming F.NoCarry condition is set (see goff for more details)")

	b.PushRegister(regY)
	b.MOVQ("res+0(FP)", regX, "dereference res")

	generateReduceASM(b, F, regT, regX)
	return nil
}

func generateDoubleASM(b *bavard.Assembly, F *field) error {
	// reset register state
	b.Reset()

	b.FuncHeader("double"+F.ElementName, 16)

	regX := b.PopRegister()
	regY := b.PopRegister()
	b.MOVQ("res+0(FP)", regX, "dereference x")
	b.MOVQ("y+8(FP)", regY, "dereference y")

	// t = y
	regT := make([]bavard.Register, F.NbWords)
	for i := 0; i < F.NbWords; i++ {
		regT[i] = b.PopRegister()
		b.MOVQ(regY.At(i), regT[i], fmt.Sprintf("t[%d] = y[%d]", i, i))
	}

	// t = t + t
	b.ADDQ(regT[0], regT[0])
	for i := 1; i < F.NbWords; i++ {
		b.ADCQ(regT[i], regT[i])
	}

	b.Comment("note that we don't check for the carry here, as this code was generated assuming F.NoCarry condition is set (see goff for more details)")

	b.PushRegister(regY)

	generateReduceASM(b, F, regT, regX)

	return nil
}
