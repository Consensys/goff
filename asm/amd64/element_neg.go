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

func (f *FFAmd64) generateNeg() {
	registers := FnHeader("neg", 0, 16)

	// labels
	nonZero := NewLabel()

	// registers
	x := registers.Pop()
	r := registers.Pop()
	q := registers.Pop()
	t := registers.PopN(f.NbWords)

	MOVQ("res+0(FP)", r)
	MOVQ("x+8(FP)", x)

	// t = x
	f.Mov(x, t)

	// x = t[0] | ... | t[n]
	MOVQ(t[0], x)
	for i := 1; i < f.NbWords; i++ {
		ORQ(t[i], x)
	}

	TESTQ(x, x)

	// if x != 0, we jump to nonzero label
	JNE(nonZero)
	// if x == 0, we set the result to zero and return
	for i := 0; i < f.NbWords/2; i++ {
		MOVQ(x, r.At(i))
	}
	RET()

	LABEL(nonZero)

	// z = x - q
	for i := 0; i < f.NbWords; i++ {
		MOVQ(f.Q[i], q)
		if i == 0 {
			SUBQ(t[i], q)
		} else {
			SBBQ(t[i], q)
		}
		MOVQ(q, r.At(i))
	}

	RET()

}
