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

func (f *FFAmd64) generateFromMont() {
	stackSize := 8
	if f.NbWords > SmallModulus {
		stackSize = f.NbWords * 8
	}
	registers := f.FnHeader("fromMont", stackSize, 8, amd64.DX, amd64.AX)
	f.WriteLn("NO_LOCAL_POINTERS")
	f.WriteLn(`
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

	noAdx := f.NewLabel()
	// check ADX instruction support
	f.CMPB("·supportAdx(SB)", 1)
	f.JNE(noAdx)

	// registers
	t := registers.PopN(f.NbWords)
	r := registers.Pop()

	f.MOVQ("res+0(FP)", r)

	// 	for i=0 to N-1
	//     t[i] = a[i]
	f.Mov(r, t)

	var tmp amd64.Register
	hasRegisters := registers.Available() > 0
	if !hasRegisters {
		tmp = r
	} else {
		tmp = registers.Pop()
	}
	for i := 0; i < f.NbWords; i++ {

		f.XORQ(amd64.DX, amd64.DX)

		// m := t[0]*q'[0] mod W
		regM := amd64.DX
		f.MOVQ(t[0], amd64.DX)
		f.MULXQ(f.qInv0(), regM, amd64.AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		f.XORQ(amd64.AX, amd64.AX)

		// C,_ := t[0] + m*q[0]
		f.Comment("C,_ := t[0] + m*q[0]")

		f.MULXQ(f.qAt(0), amd64.AX, tmp)
		f.ADCXQ(t[0], amd64.AX)
		f.MOVQ(tmp, t[0])

		f.Comment("for j=1 to N-1")
		f.Comment("    (C,t[j-1]) := t[j] + m*q[j] + C")

		// for j=1 to N-1
		//    (C,t[j-1]) := t[j] + m*q[j] + C
		for j := 1; j < f.NbWords; j++ {
			f.ADCXQ(t[j], t[j-1])
			f.MULXQ(f.qAt(j), amd64.AX, t[j])
			f.ADOXQ(amd64.AX, t[j-1])
		}
		f.MOVQ(0, amd64.AX)
		f.ADCXQ(amd64.AX, t[f.NbWordsLastIndex])
		f.ADOXQ(amd64.AX, t[f.NbWordsLastIndex])

	}

	if !hasRegisters {
		f.MOVQ("res+0(FP)", r)
	} else {
		registers.Push(tmp)
	}
	// ---------------------------------------------------------------------------------------------
	// reduce
	f.Reduce(&registers, t, r)
	f.RET()

	// No adx
	f.LABEL(noAdx)
	f.MOVQ("res+0(FP)", amd64.AX)
	f.MOVQ(amd64.AX, "(SP)")
	f.WriteLn("CALL ·_fromMontGeneric(SB)")
	f.RET()

}
