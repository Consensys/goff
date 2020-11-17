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

func (f *FFAmd64) generateSub() {
	registers := FnHeader("sub", 0, 24)

	// registers
	t := registers.PopN(f.NbWords)
	x := registers.Pop()
	y := registers.Pop()

	MOVQ("x+8(FP)", x)
	f.Mov(x, t)

	// z = x - y mod q
	MOVQ("y+16(FP)", y)
	f.Sub(y, t)

	if f.NbWords > 6 {
		f.ReduceAfterSub(&registers, t, false)
	} else {
		f.ReduceAfterSub(&registers, t, true)
	}

	r := registers.Pop()
	MOVQ("res+0(FP)", r)
	f.Mov(t, r)

	RET()

}

func (f *FFAmd64) ReduceAfterSub(registers *Registers, t []Register, noJump bool) {
	if noJump {
		q := registers.PopN(f.NbWords)
		r := registers.Pop()
		f.Mov(f.Q, q)
		MOVQ(0, r)
		// overwrite with 0 if borrow is set
		for i := 0; i < f.NbWords; i++ {
			CMOVQCC(r, q[i])
		}

		// add registers (q or 0) to t, and set to result
		f.Add(q, t)

		registers.Push(r)
		registers.Push(q...)
	} else {
		noReduce := NewLabel()

		JCC(noReduce)
		for i := 0; i < f.NbWords; i++ {
			if i == 0 {
				ADDQ(f.qAt(i), t[i])
			} else {
				ADCQ(f.qAt(i), t[i])
			}
		}
		LABEL(noReduce)

	}
}
