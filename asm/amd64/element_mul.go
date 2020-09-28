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

import . "github.com/consensys/bavard/amd64"

func (f *FFAmd64) MulADX(registers *Registers, yat, xat func(int) string, uglyHook func(int)) []Register {
	// registers
	t := registers.PopN(f.NbWords)
	A := registers.Pop()

	for i := 0; i < f.NbWords; i++ {
		XORQ(DX, DX)

		MOVQ(yat(i), DX)
		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < f.NbWords; j++ {
			xj := xat(j)

			reg := A
			if i == 0 {
				if j == 0 {
					MULXQ(xj, t[j], t[j+1])
				} else if j != f.NbWordsLastIndex {
					reg = t[j+1]
				}
			} else if j != 0 {
				ADCXQ(A, t[j])
			}

			if !(i == 0 && j == 0) {
				MULXQ(xj, AX, reg)
				ADOXQ(AX, t[j])
			}
		}
		if uglyHook != nil {
			uglyHook(i)
		}

		Comment("add the last carries to " + string(A))
		MOVQ(0, DX)
		ADCXQ(DX, A)
		ADOXQ(DX, A)

		// m := t[0]*q'[0] mod W
		m := DX
		MOVQ(t[0], DX)
		MULXQ(f.qInv0(), m, AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		XORQ(AX, AX)

		// C,_ := t[0] + m*q[0]
		Comment("C,_ := t[0] + m*q[0]")

		needPop := false
		if registers.Available() == 0 {
			needPop = true
			PUSHQ(A)
			registers.Push(A)
		}
		tmp := registers.Pop()
		MULXQ(f.qAt(0), AX, tmp)
		ADCXQ(t[0], AX)
		MOVQ(tmp, t[0])
		registers.Push(tmp)
		if needPop {
			A = registers.Pop()
			POPQ(A)
		}

		Comment("for j=1 to N-1")
		Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < f.NbWords; j++ {
			ADCXQ(t[j], t[j-1])
			MULXQ(f.qAt(j), AX, t[j])
			ADOXQ(AX, t[j-1])
		}
		MOVQ(0, AX)
		ADCXQ(AX, t[f.NbWordsLastIndex])
		ADOXQ(A, t[f.NbWordsLastIndex])
	}

	// free registers
	registers.Push(A)

	return t
}

func (f *FFAmd64) generateInnerMul(registers *Registers, isSquare bool) {

	noAdx := NewLabel()

	// check ADX instruction support
	CMPB("·supportAdx(SB)", 1)
	JNE(noAdx)
	{
		var t []Register
		if isSquare {
			x := registers.Pop()
			MOVQ("x+8(FP)", x)

			xat := func(i int) string {
				return x.At(i)
			}
			t = f.MulADX(registers, xat, xat, nil)
			registers.Push(x, x)
		} else {
			x := registers.Pop()
			y := registers.Pop()
			MOVQ("x+8(FP)", x)
			MOVQ("y+16(FP)", y)

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
		MOVQ("res+0(FP)", r)
		f.Reduce(registers, t, r)
		RET()
		registers.Push(r)
	}

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		LABEL(noAdx)
		registers := NewRegisters()
		registers.Remove(AX)
		registers.Remove(DX)
		x := registers.Pop()
		y := registers.Pop()
		if isSquare {
			MOVQ("x+8(FP)", x)
			MOVQ("x+8(FP)", y)
		} else {
			MOVQ("x+8(FP)", x)
			MOVQ("y+16(FP)", y)
		}

		f.mulNoAdx(&registers, x, y)
	}
}

func (f *FFAmd64) generateMul() {
	stackSize := 0
	if f.NbWords > SmallModulus {
		stackSize = f.NbWords * 8
	}
	registers := FnHeader("mul", stackSize, 24, DX, AX)
	WriteLn(`
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
		f.generateInnerMulLarge(&registers, false)
	} else {
		f.generateInnerMul(&registers, false)
	}

}

func (f *FFAmd64) generateInnerMulLarge(registers *Registers, isSquare bool) {
	WriteLn("NO_LOCAL_POINTERS")
	noAdx := NewLabel()
	// check ADX instruction support
	CMPB("·supportAdx(SB)", 1)
	JNE(noAdx)

	// registers
	t := registers.PopN(f.NbWords)
	A := registers.Pop()

	for i := 0; i < f.NbWords; i++ {

		XORQ(DX, DX)
		yi := DX
		if isSquare {
			MOVQ("x+8(FP)", yi)
		} else {
			MOVQ("y+16(FP)", yi)
		}
		MOVQ(yi.At(i), yi)
		// for j=0 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		for j := 0; j < f.NbWords; j++ {
			xj := AX
			MOVQ("x+8(FP)", xj)
			MOVQ(xj.At(j), xj)

			reg := A
			if i == 0 {
				if j == 0 {
					MULXQ(xj, t[j], t[j+1])
				} else if j != f.NbWordsLastIndex {
					reg = t[j+1]
				}
			} else if j != 0 {
				ADCXQ(A, t[j])
			}

			if !(i == 0 && j == 0) {
				MULXQ(xj, AX, reg)
				ADOXQ(AX, t[j])
			}
		}

		Comment("add the last carries to " + string(A))
		MOVQ(0, DX)
		ADCXQ(DX, A)
		ADOXQ(DX, A)
		PUSHQ(A)

		// m := t[0]*q'[0] mod W
		regM := DX
		MOVQ(t[0], DX)
		MULXQ(f.qInv0(), regM, AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		XORQ(AX, AX)

		// C,_ := t[0] + m*q[0]
		Comment("C,_ := t[0] + m*q[0]")
		MULXQ(f.qAt(0), AX, A)
		ADCXQ(t[0], AX)
		MOVQ(A, t[0])

		Comment("for j=1 to N-1")
		Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < f.NbWords; j++ {
			ADCXQ(t[j], t[j-1])
			MULXQ(f.qAt(j), AX, t[j])
			ADOXQ(AX, t[j-1])
		}

		POPQ(A)
		MOVQ(0, AX)
		ADCXQ(AX, t[f.NbWordsLastIndex])
		ADOXQ(A, t[f.NbWordsLastIndex])
	}

	// free registers
	registers.Push(A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := registers.Pop()
	MOVQ("res+0(FP)", r)
	f.reduceLarge(t, r)
	RET()

	// No adx
	LABEL(noAdx)
	MOVQ("res+0(FP)", AX)
	MOVQ(AX, "(SP)")
	MOVQ("x+8(FP)", AX)
	MOVQ(AX, "8(SP)")
	if isSquare {
		WriteLn("CALL ·_squareGeneric(SB)")
		RET()
	} else {
		MOVQ("y+16(FP)", AX)
		MOVQ(AX, "16(SP)")
		WriteLn("CALL ·_mulGeneric(SB)")
		RET()
	}

}

func (f *FFAmd64) mulNoAdx(registers *Registers, x, y Register) {
	// registers
	t := registers.PopN(f.NbWords)
	C := registers.Pop()
	yi := registers.Pop()
	A := registers.Pop()
	m := registers.Pop()

	for i := 0; i < f.NbWords; i++ {
		// (A,t[0]) := t[0] + x[0]*y[{{$i}}]
		MOVQ(x.At(0), AX)
		MOVQ(y.At(i), yi)
		MULQ(yi)
		if i != 0 {
			ADDQ(AX, t[0])
			ADCQ(0, DX)
		} else {
			MOVQ(AX, t[0])
		}
		MOVQ(DX, A)

		// m := t[0]*q'[0] mod W
		MOVQ(f.qInv0(), m)
		IMULQ(t[0], m)

		// C,_ := t[0] + m*q[0]
		MOVQ(f.Q[0], AX)
		MULQ(m)
		ADDQ(t[0], AX)
		ADCQ(0, DX)
		MOVQ(DX, C)

		// for j=1 to N-1
		//    (A,t[j])  := t[j] + x[j]*y[i] + A
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < f.NbWords; j++ {
			MOVQ(x.At(j), AX)
			MULQ(yi)
			if i != 0 {
				ADDQ(A, t[j])
				ADCQ(0, DX)
				ADDQ(AX, t[j])
				ADCQ(0, DX)
			} else {
				MOVQ(A, t[j])
				ADDQ(AX, t[j])
				ADCQ(0, DX)
			}
			MOVQ(DX, A)

			MOVQ(f.Q[j], AX)
			MULQ(m)
			ADDQ(t[j], C)
			ADCQ(0, DX)
			ADDQ(AX, C)
			ADCQ(0, DX)
			MOVQ(C, t[j-1])
			MOVQ(DX, C)
		}

		ADDQ(C, A)
		MOVQ(A, t[f.NbWordsLastIndex])

	}

	// ---------------------------------------------------------------------------------------------
	// reduce
	registers.Push(C, yi, A, m, y)

	MOVQ("res+0(FP)", x)
	f.Reduce(registers, t, x)
	RET()
}
