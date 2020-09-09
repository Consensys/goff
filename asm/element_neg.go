// Copyright 2020 ConsenSys AG
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

package asm

func generateNeg() {
	fnHeader("neg", 0, 16)

	// labels
	nonZero := newLabel()

	// registers
	x := popRegister()
	r := popRegister()
	q := popRegister()
	t := popRegisters(nbWords)

	movq("res+0(FP)", r)
	movq("x+8(FP)", x)

	// t = x
	_mov(x, t)

	// x = t[0] | ... | t[n]
	movq(t[0], x)
	for i := 1; i < nbWords; i++ {
		orq(t[i], x)
	}

	testq(x, x)

	// if x != 0, we jump to nonzero label
	jne(nonZero)
	// if x == 0, we set the result to zero and return
	for i := 0; i < nbWords/2; i++ {
		movq(x, r.at(i))
	}
	ret()

	label(nonZero)

	// z = x - q
	for i := 0; i < nbWords; i++ {
		movq(modulus[i], q)
		if i == 0 {
			subq(t[i], q)
		} else {
			sbbq(t[i], q)
		}
		movq(q, r.at(i))
	}

	ret()

}
