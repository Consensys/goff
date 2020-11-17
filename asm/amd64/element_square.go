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

import "github.com/consensys/bavard/amd64"

func (f *FFAmd64) generateSquare() {
	stackSize := 0
	if f.NbWords > SmallModulus {
		stackSize = f.NbWords * 8
	}
	registers := f.FnHeader("square", stackSize, 16, amd64.DX, amd64.AX)
	f.WriteLn(`
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

	noAdx := f.NewLabel()
	// check ADX instruction support
	f.CMPB("·supportAdx(SB)", 1)
	f.JNE(noAdx)

	// registers
	t := registers.PopN(f.NbWords)
	x := registers.Pop()
	A := registers.Pop()

	f.MOVQ("x+8(FP)", x)
	for i := 0; i < f.NbWords; i++ {

		f.XORQ(amd64.AX, amd64.AX)

		f.MOVQ(x.At(i), amd64.DX)

		// instead of
		// for j=i+1 to N-1
		//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
		// we first add the x[j]*x[i] to a temporary u (set of registers)
		// set double it, before doing
		// for j=i+1 to N-1
		//     A,t[j] = u[j] + t[j] + A
		if i != f.NbWordsLastIndex {
			u := make([]amd64.Register, (f.NbWords - i - 1))
			for i := 0; i < len(u); i++ {
				u[i] = registers.Pop()
			}
			offset := i + 1

			// 1- compute u = x[j] * x[i]
			// for j=i+1 to N-1
			//     A,u[j] = x[j]*x[i] + A
			if (i + 1) == f.NbWordsLastIndex {
				f.MULXQ(x.At(i+1), u[0], A)
			} else {
				for j := i + 1; j < f.NbWords; j++ {
					yj := x.At(j)
					if j == i+1 {
						// first iteration
						f.MULXQ(yj, u[j-offset], u[j+1-offset])
					} else {
						if j == f.NbWordsLastIndex {
							f.MULXQ(yj, amd64.AX, A)
						} else {
							f.MULXQ(yj, amd64.AX, u[j+1-offset])
						}
						f.ADCXQ(amd64.AX, u[j-offset])
					}
				}
				f.MOVQ(0, amd64.AX)
				f.ADCXQ(amd64.AX, A)
				f.XORQ(amd64.AX, amd64.AX)
			}

			if i == 0 {
				// C, t[i] = x[i] * x[i] + t[i]
				f.MULXQ(amd64.DX, t[i], amd64.DX)

				// when i == 0, T is not set yet
				// so  we can use f.ADOXQ carry chain to propagate C from x[i] * x[i] + t[i] (dx)

				// for j=i+1 to N-1
				// 		C, t[j] = u[j] + u[j] + t[j] + C
				for j := 0; j < len(u); j++ {
					f.ADCXQ(u[j], u[j])
					f.MOVQ(u[j], t[j+offset])
					if j == 0 {
						f.ADOXQ(amd64.DX, t[j+offset])
					} else {
						f.ADOXQ(amd64.AX, t[j+offset])
					}
				}

				f.ADCXQ(A, A)
				f.ADOXQ(amd64.AX, A)

			} else {
				// i != 0 so T is set.
				// we first use f.ADOXQ carry chain to perform t = u + u + t
				for j := 0; j < len(u); j++ {
					f.ADCXQ(u[j], u[j])
					f.ADOXQ(u[j], t[j+offset])
				}

				f.ADCXQ(A, A)
				f.ADOXQ(amd64.AX, A)

				// reset flags
				f.XORQ(amd64.AX, amd64.AX)

				// C, t[i] = x[i] * x[i] + t[i]
				f.MULXQ(amd64.DX, amd64.AX, amd64.DX)
				f.ADOXQ(amd64.AX, t[i])
				f.MOVQ(0, amd64.AX)

				// propagate C
				for j := i + 1; j < f.NbWords; j++ {
					if j == i+1 {
						f.ADOXQ(amd64.DX, t[j])
					} else {
						f.ADOXQ(amd64.AX, t[j])
					}
				}

				f.ADOXQ(amd64.AX, A)
			}

			registers.Push(u...)

		} else {
			// i == last index
			f.MULXQ(amd64.DX, amd64.AX, A)
			f.ADCXQ(amd64.AX, t[i])
			f.MOVQ(0, amd64.AX)
			f.ADCXQ(amd64.AX, A)
		}

		tmp := registers.Pop()
		// m := t[0]*q'[0] mod W
		regM := amd64.DX
		f.MOVQ(t[0], amd64.DX)
		f.MULXQ(f.qInv0(), regM, amd64.AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		f.XORQ(amd64.AX, amd64.AX)

		// C,_ := t[0] + m*q[0]
		f.MULXQ(f.qAt(0), amd64.AX, tmp)
		f.ADCXQ(t[0], amd64.AX)
		f.MOVQ(tmp, t[0])

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < f.NbWords; j++ {
			f.ADCXQ(t[j], t[j-1])
			f.MULXQ(f.qAt(j), amd64.AX, t[j])
			f.ADOXQ(amd64.AX, t[j-1])
		}
		f.MOVQ(0, amd64.AX)
		f.ADCXQ(amd64.AX, t[f.NbWordsLastIndex])
		f.ADOXQ(A, t[f.NbWordsLastIndex])

		registers.Push(tmp)
	}

	// free registers
	registers.Push(x, A)

	// ---------------------------------------------------------------------------------------------
	// reduce
	r := registers.Pop()
	f.MOVQ("res+0(FP)", r)
	f.Reduce(&registers, t, r)
	f.RET()

	// ---------------------------------------------------------------------------------------------
	// no MULX, ADX instructions
	{
		f.LABEL(noAdx)
		registers = amd64.NewRegisters()
		registers.Remove(amd64.AX)
		registers.Remove(amd64.DX)
		x := registers.Pop()
		y := registers.Pop()
		f.MOVQ("x+8(FP)", x)
		f.MOVQ("x+8(FP)", y)
		f.mulNoAdx(&registers, x, y)
	}

}
