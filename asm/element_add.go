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

func generateAdd() {
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("add", stackSize, 24)

	// registers
	x := popRegister()
	y := popRegister()
	r := popRegister()
	t := popRegisters(nbWords)

	movq("x+8(FP)", x)

	// t = x
	_mov(x, t)

	movq("y+16(FP)", y)

	// t = t + y = x + y
	_add(y, t)

	// dereference res
	movq("res+0(FP)", r)

	// reduce t into res
	_reduce(t, r)

	ret()

}
