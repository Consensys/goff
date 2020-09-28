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

func (f *FFAmd64) generateSquare() {
	stackSize := 0
	if f.NbWords > SmallModulus {
		stackSize = f.NbWords * 8
	}
	registers := FnHeader("square", stackSize, 16, DX, AX)
	WriteLn(`
	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// for i=0 to N-1
	// A, t[i] = x[i] * x[i] + t[i]
	// p = 0
	// for j=i+1 to N-1
	//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
	// m = t[0] * q'[0]
	// C, _ = t[0] + q[0]*m
	// for j=1 to N-1
	//     C, t[j-1] = q[j]*m +  t[j] + C
	// t[N-1] = C + A

	`)
	if f.NbWords > 6 {
		f.generateInnerMulLarge(&registers, true)
		return
	}
	if !f.NoCarrySquare {
		f.generateInnerMul(&registers, true)
		return
	}

	noAdx := NewLabel()
	// check ADX instruction support
	CMPB("·supportAdx(SB)", 1)
	JNE(noAdx)

	// registers
	t := registers.PopN(f.NbWords)
	x := registers.Pop()
	A := registers.Pop()

	MOVQ("x+8(FP)", x)
	for i := 0; i < f.NbWords; i++ {

		XORQ(AX, AX)

		MOVQ(x.At(i), DX)

		// instead of
		// for j=i+1 to N-1
		//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
		// we first add the x[j]*x[i] to a temporary u (set of registers)
		// set double it, before doing
		// for j=i+1 to N-1
		//     A,t[j] = u[j] + t[j] + A
		if i != f.NbWordsLastIndex {
			u := make([]Register, (f.NbWords - i - 1))
			for i := 0; i < len(u); i++ {
				u[i] = registers.Pop()
			}
			offset := i + 1

			// 1- compute u = x[j] * x[i]
			// for j=i+1 to N-1
			//     A,u[j] = x[j]*x[i] + A
			if (i + 1) == f.NbWordsLastIndex {
				MULXQ(x.At(i+1), u[0], A)
			} else {
				for j := i + 1; j < f.NbWords; j++ {
					yj := x.At(j)
					if j == i+1 {
						// first iteration
						MULXQ(yj, u[j-offset], u[j+1-offset])
					} else {
						if j == f.NbWordsLastIndex {
							MULXQ(yj, AX, A)
						} else {
							MULXQ(yj, AX, u[j+1-offset])
						}
						ADCXQ(AX, u[j-offset])
					}
				}
				MOVQ(0, AX)
				ADCXQ(AX, A)
				XORQ(AX, AX)
			}

			if i == 0 {
				// C, t[i] = x[i] * x[i] + t[i]
				MULXQ(DX, t[i], DX)

				// when i == 0, T is not set yet
				// so  we can use ADOXQ carry chain to propagate C from x[i] * x[i] + t[i] (dx)

				// for j=i+1 to N-1
				// 		C, t[j] = u[j] + u[j] + t[j] + C
				for j := 0; j < len(u); j++ {
					ADCXQ(u[j], u[j])
					MOVQ(u[j], t[j+offset])
					if j == 0 {
						ADOXQ(DX, t[j+offset])
					} else {
						ADOXQ(AX, t[j+offset])
					}
				}

				ADCXQ(A, A)
				ADOXQ(AX, A)

			} else {
				// i != 0 so T is set.
				// we first use ADOXQ carry chain to perform t = u + u + t
				for j := 0; j < len(u); j++ {
					ADCXQ(u[j], u[j])
					ADOXQ(u[j], t[j+offset])
				}

				ADCXQ(A, A)
				ADOXQ(AX, A)

				// reset flags
				XORQ(AX, AX)

				// C, t[i] = x[i] * x[i] + t[i]
				MULXQ(DX, AX, DX)
				ADOXQ(AX, t[i])
				MOVQ(0, AX)

				// propagate C
				for j := i + 1; j < f.NbWords; j++ {
					if j == i+1 {
						ADOXQ(DX, t[j])
					} else {
						ADOXQ(AX, t[j])
					}
				}

				ADOXQ(AX, A)
			}

			registers.Push(u...)

		} else {
			// i == last index
			MULXQ(DX, AX, A)
			ADCXQ(AX, t[i])
			MOVQ(0, AX)
			ADCXQ(AX, A)
		}

		tmp := registers.Pop()
		// m := t[0]*q'[0] mod W
		regM := DX
		MOVQ(t[0], DX)
		MULXQ(f.qInv0(), regM, AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		XORQ(AX, AX)

		// C,_ := t[0] + m*q[0]
		MULXQ(f.qAt(0), AX, tmp)
		ADCXQ(t[0], AX)
		MOVQ(tmp, t[0])

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

		registers.Push(tmp)
	}

	// free registers
	registers.Push(x, A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := registers.Pop()
	MOVQ("res+0(FP)", r)
	f.Reduce(&registers, t, r)
	RET()

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		LABEL(noAdx)
		registers = NewRegisters()
		registers.Remove(AX)
		registers.Remove(DX)
		x := registers.Pop()
		y := registers.Pop()
		MOVQ("x+8(FP)", x)
		MOVQ("x+8(FP)", y)
		f.mulNoAdx(&registers, x, y)
	}

}
