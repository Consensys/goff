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

func (f *ffAmd64) generateFromMont() {
	stackSize := 8
	if f.NbWords > SmallModulus {
		stackSize = f.NbWords * 8
	}
	registers := FnHeader("fromMont", stackSize, 8, DX, AX)
	WriteLn("NO_LOCAL_POINTERS")
	WriteLn(`
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

	noAdx := NewLabel()
	// check ADX instruction support
	CMPB("·supportAdx(SB)", 1)
	JNE(noAdx)

	// registers
	t := registers.PopN(f.NbWords)
	r := registers.Pop()

	MOVQ("res+0(FP)", r)

	// 	for i=0 to N-1
	//     t[i] = a[i]
	f.Mov(r, t)

	var tmp Register
	hasRegisters := registers.Available() > 0
	if !hasRegisters {
		tmp = r
	} else {
		tmp = registers.Pop()
	}
	for i := 0; i < f.NbWords; i++ {

		XORQ(DX, DX)

		// m := t[0]*q'[0] mod W
		regM := DX
		MOVQ(t[0], DX)
		MULXQ(f.qInv0(), regM, AX, "m := t[0]*q'[0] mod W")

		// clear the carry flags
		XORQ(AX, AX)

		// C,_ := t[0] + m*q[0]
		Comment("C,_ := t[0] + m*q[0]")

		MULXQ(f.qAt(0), AX, tmp)
		ADCXQ(t[0], AX)
		MOVQ(tmp, t[0])

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
		ADOXQ(AX, t[f.NbWordsLastIndex])

	}

	if !hasRegisters {
		MOVQ("res+0(FP)", r)
	} else {
		registers.Push(tmp)
	}
	// ---------------------------------------------------------------------------------------------
	// reduce
	f.Reduce(&registers, t, r)
	RET()

	// No adx
	LABEL(noAdx)
	MOVQ("res+0(FP)", AX)
	MOVQ(AX, "(SP)")
	WriteLn("CALL ·_fromMontGeneric(SB)")
	RET()

}
