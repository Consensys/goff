// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package amd64

import (
	"fmt"

	"github.com/consensys/bavard/amd64"
)

func (f *FFAmd64) MulADX(registers *amd64.Registers, yat, xat func(int) string, uglyHook func(int)) []amd64.Register {
	// registers
	t := registers.PopN(f.NbWords)
	A := registers.Pop()

	f.LabelRegisters("t", t...)

	for i := 0; i < f.NbWords; i++ {
		f.Comment("clear the flags")
		f.XORQ(amd64.AX, amd64.AX)

		if yat == nil {
			f.POPQ(amd64.DX)
		} else {
			f.MOVQ(yat(i), amd64.DX)
		}

		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < f.NbWords; j++ {
			f.Comment(fmt.Sprintf("(A,t[%[1]d])  := t[%[1]d] + x[%[1]d]*y[%[2]d] + A", j, i))
			xj := xat(j)

			reg := A
			if i == 0 {
				if j == 0 {
					f.MULXQ(xj, t[j], t[j+1])
				} else if j != f.NbWordsLastIndex {
					reg = t[j+1]
				}
			} else if j != 0 {
				f.ADCXQ(A, t[j])
			}

			if !(i == 0 && j == 0) {
				f.MULXQ(xj, amd64.AX, reg)
				f.ADOXQ(amd64.AX, t[j])
			}
		}
		if uglyHook != nil {
			uglyHook(i)
		}

		f.Comment("A += carries from ADCXQ and ADOXQ")
		f.MOVQ(0, amd64.AX)
		if i != 0 {
			f.ADCXQ(amd64.AX, A)
		}
		f.ADOXQ(amd64.AX, A)

		// m := t[0]*q'[0] mod W
		f.Comment("m := t[0]*q'[0] mod W")
		m := amd64.DX
		// f.MOVQ(t[0], m)
		// f.MULXQ(f.qInv0(), m, amd64.AX)
		f.MOVQ(f.qInv0(), m)
		f.IMULQ(t[0], m)

		// clear the carry flags
		f.Comment("clear the flags")
		f.XORQ(amd64.AX, amd64.AX)

		// C,_ := t[0] + m*q[0]
		f.Comment("C,_ := t[0] + m*q[0]")

		needPop := false
		if registers.Available() == 0 {
			needPop = true
			f.PUSHQ(A)
			registers.Push(A)
		}
		tmp := registers.Pop()
		f.MULXQ(f.qAt(0), amd64.AX, tmp)
		f.ADCXQ(t[0], amd64.AX)
		f.MOVQ(tmp, t[0])
		registers.Push(tmp)
		if needPop {
			A = registers.Pop()
			f.POPQ(A)
		}

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < f.NbWords; j++ {
			f.Comment(fmt.Sprintf("(C,t[%[1]d]) := t[%[2]d] + m*q[%[2]d] + C", j-1, j))
			f.ADCXQ(t[j], t[j-1])
			f.MULXQ(f.qAt(j), amd64.AX, t[j])
			f.ADOXQ(amd64.AX, t[j-1])
		}

		f.Comment(fmt.Sprintf("t[%d] = C + A", f.NbWordsLastIndex))
		f.MOVQ(0, amd64.AX)
		f.ADCXQ(amd64.AX, t[f.NbWordsLastIndex])
		f.ADOXQ(A, t[f.NbWordsLastIndex])
	}

	// free registers
	registers.Push(A)

	return t
}

func (f *FFAmd64) generateInnerMul(registers *amd64.Registers, isSquare bool) {
	f.WriteLn("NO_LOCAL_POINTERS")
	noAdx := f.NewLabel()
	// check ADX instruction support
	f.CMPB("·supportAdx(SB)", 1)
	f.JNE(noAdx)
	{
		var t []amd64.Register
		if isSquare {
			x := registers.Pop()
			f.MOVQ("x+8(FP)", x)

			xat := func(i int) string {
				return x.At(i)
			}
			t = f.MulADX(registers, xat, xat, nil)
			registers.Push(x, x)
		} else {
			x := registers.Pop()
			y := registers.Pop()
			f.MOVQ("x+8(FP)", x)
			f.MOVQ("y+16(FP)", y)

			yat := func(i int) string {
				return y.At(i)
			}
			xat := func(i int) string {
				return x.At(i)
			}
			t = f.MulADX(registers, yat, xat, nil)
			registers.Push(x, y)
		}

		r := registers.Pop()
		// ---------------------------------------------------------------------------------------------
		// reduce
		f.MOVQ("res+0(FP)", r)
		f.Reduce(registers, t, r)
		f.RET()
		registers.Push(r)
	}

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		f.LABEL(noAdx)

		f.MOVQ("res+0(FP)", amd64.AX)
		f.MOVQ(amd64.AX, "(SP)")
		f.MOVQ("x+8(FP)", amd64.AX)
		f.MOVQ(amd64.AX, "8(SP)")
		f.MOVQ("y+16(FP)", amd64.AX)
		f.MOVQ(amd64.AX, "16(SP)")
		f.WriteLn("CALL ·_mulGeneric(SB)")
		f.RET()

	}
}

func (f *FFAmd64) generateMul() {
	stackSize := 3 * 8
	if f.NbWords > SmallModulus {
		stackSize = f.NbWords * 8
	}
	f.Comment("mul(res, x, y *Element)")
	registers := f.FnHeader("mul", stackSize, 24, amd64.DX, amd64.AX)
	f.WriteLn(`
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
	if f.NbWords > SmallModulus {
		f.generateInnerMulLarge(&registers)
	} else {
		f.generateInnerMul(&registers, false)
	}

}

func (f *FFAmd64) generateInnerMulLarge(registers *amd64.Registers) {
	f.WriteLn("NO_LOCAL_POINTERS")
	noAdx := f.NewLabel()
	// check ADX instruction support
	f.CMPB("·supportAdx(SB)", 1)
	f.JNE(noAdx)

	// registers
	t := registers.PopN(f.NbWords)
	A := registers.Pop()
	f.LabelRegisters("A", A)
	f.LabelRegisters("t", t...)

	x := make([]amd64.Register, f.NbWords)
	for i := 0; i < f.NbWords; i++ {
		// use stack
		x[i] = amd64.Register(fmt.Sprintf("x%d-%d(SP)", i, 8+i*8))

	}

	for i := 0; i < f.NbWords; i++ {
		f.Comment("clear the flags")
		f.XORQ(amd64.DX, amd64.DX)
		yi := amd64.DX
		f.MOVQ("y+16(FP)", yi)
		f.MOVQ(yi.At(i), yi)

		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < f.NbWords; j++ {
			f.Comment(fmt.Sprintf("(A,t[%[1]d])  := t[%[1]d] + x[%[1]d]*y[%[2]d] + A", j, i))
			xj := amd64.AX
			if i == 0 {
				xj = A
				if j == 0 {
					f.Comment(fmt.Sprintf("using A(%s) to store x", A))
					f.MOVQ("x+8(FP)", xj)
				}
				f.MOVQ(xj.At(j), amd64.AX)
				f.MOVQ(amd64.AX, x[j])
			}

			reg := A
			if i == 0 && j == 0 {
				f.MULXQ(amd64.AX, t[j], t[j+1])
			}
			if i != 0 && j != 0 {
				f.ADCXQ(A, t[j])
			}
			if i != 0 || j != 0 {
				if i == 0 && j != f.NbWordsLastIndex {
					reg = t[j+1]
				}
				if i == 0 {
					f.MULXQ(amd64.AX, amd64.AX, reg)
				} else {
					f.MULXQ(x[j], amd64.AX, reg)
				}
				f.ADOXQ(amd64.AX, t[j])
			}
		}

		f.Comment("A += carries from ADCXQ and ADOXQ")
		f.MOVQ(0, amd64.DX)
		if i != 0 {
			f.ADCXQ(amd64.DX, A)
		}
		f.ADOXQ(amd64.DX, A)
		// TODO need to avoid PUSHQ / POPQ as we need to be careful with stack usage of caller funcs, and usage of virtual (SP)
		f.PUSHQ(A)

		// m := t[0]*q'[0] mod W
		f.Comment("m := t[0]*q'[0] mod W")
		m := amd64.DX
		// f.MOVQ(t[0], amd64.DX)
		// f.MULXQ(f.qInv0(), m, amd64.AX)
		f.MOVQ(f.qInv0(), m)
		f.IMULQ(t[0], m)

		// clear the carry flags
		f.Comment("clear the flags")
		f.XORQ(amd64.AX, amd64.AX)

		// C,_ := t[0] + m*q[0]
		f.Comment("C,_ := t[0] + m*q[0]")
		f.MULXQ(f.qAt(0), amd64.AX, A)
		f.ADCXQ(t[0], amd64.AX)
		f.MOVQ(A, t[0])

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < f.NbWords; j++ {
			f.Comment(fmt.Sprintf("(C,t[%[1]d]) := t[%[2]d] + m*q[%[2]d] + C", j-1, j))
			f.ADCXQ(t[j], t[j-1])
			f.MULXQ(f.qAt(j), amd64.AX, t[j])
			f.ADOXQ(amd64.AX, t[j-1])
		}

		f.Comment(fmt.Sprintf("t[%d] = C + A", f.NbWordsLastIndex))
		f.POPQ(A)
		f.MOVQ(0, amd64.AX)
		f.ADCXQ(amd64.AX, t[f.NbWordsLastIndex])
		f.ADOXQ(A, t[f.NbWordsLastIndex])
	}

	// free registers
	registers.Push(A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := registers.Pop()
	f.MOVQ("res+0(FP)", r)
	f.reduceLarge(t, r)
	f.RET()

	// No adx
	f.LABEL(noAdx)
	f.MOVQ("res+0(FP)", amd64.AX)
	f.MOVQ(amd64.AX, "(SP)")
	f.MOVQ("x+8(FP)", amd64.AX)
	f.MOVQ(amd64.AX, "8(SP)")
	f.MOVQ("y+16(FP)", amd64.AX)
	f.MOVQ(amd64.AX, "16(SP)")
	f.WriteLn("CALL ·_mulGeneric(SB)")
	f.RET()

}
