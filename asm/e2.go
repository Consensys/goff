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

import "strings"

func generateAddE2() {
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("add"+strings.ToUpper(elementName), stackSize, 24)

	// registers
	x := popRegister()
	y := popRegister()
	r := popRegister()
	t := popRegisters(nbWords)

	movq("x+8(FP)", x)

	// move t = x
	_mov(x, t)

	movq("y+16(FP)", y)

	// t = t + y = x + y
	_add(y, t)

	// reduce
	movq("res+0(FP)", r)
	_reduce(t, r)

	// move x+offset(nbWords) into t
	_mov(x, t, nbWords)

	// add y+offset(nbWords) into t
	_add(y, t, nbWords)

	// reduce t into r with offset nbWords
	_reduce(t, r, nbWords)

	ret()

}

func generateDoubleE2() {
	// func header
	stackSize := 0
	if nbWords > smallModulus {
		stackSize = nbWords * 8
	}
	fnHeader("double"+strings.ToUpper(elementName), stackSize, 16)

	// registers
	x := popRegister()
	r := popRegister()
	t := popRegisters(nbWords)

	movq("res+0(FP)", r)
	movq("x+8(FP)", x)

	_mov(x, t)
	_add(t, t)
	_reduce(t, r)
	_mov(x, t, nbWords)
	_add(t, t)
	_reduce(t, r, nbWords)

	ret()
}

func generateNegE2() {
	fnHeader("neg"+strings.ToUpper(elementName), 0, 16)

	nonZeroA := newLabel()
	nonZeroB := newLabel()
	B := newLabel()

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
	jne(nonZeroA)

	// if x == 0, we set the result to zero and continue
	for i := 0; i < nbWords; i++ {
		movq(x, r.at(i+nbWords))
	}
	jmp(B)

	label(nonZeroA)

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

	label(B)
	movq("x+8(FP)", x)
	_mov(x, t, nbWords)

	// x = t[0] | ... | t[n]
	movq(t[0], x)
	for i := 1; i < nbWords; i++ {
		orq(t[i], x)
	}

	testq(x, x)

	// if x != 0, we jump to nonzero label
	jne(nonZeroB)

	// if x == 0, we set the result to zero and return
	for i := 0; i < nbWords; i++ {
		movq(x, r.at(i+nbWords))
	}
	ret()

	label(nonZeroB)

	// z = x - q
	for i := 0; i < nbWords; i++ {
		movq(modulus[i], q)
		if i == 0 {
			subq(t[i], q)
		} else {
			sbbq(t[i], q)
		}
		movq(q, r.at(i+nbWords))
	}

	ret()

}

func generateSubE2() {
	fnHeader("sub"+strings.ToUpper(elementName), 0, 24)

	// registers
	t := popRegisters(nbWords)
	x := popRegister()
	y := popRegister()

	movq("x+8(FP)", x)
	movq("y+16(FP)", y)

	_mov(x, t)

	// z = x - y mod q
	// move t = x
	_sub(y, t)

	if nbWords > 6 {
		_reduceAfterSub(t, false)
	} else {
		_reduceAfterSub(t, true)
	}

	r := popRegister()
	movq("res+0(FP)", r)
	_mov(t, r)
	pushRegister(r)

	_mov(x, t, nbWords)

	// z = x - y mod q
	// move t = x
	_sub(y, t, nbWords)

	if nbWords > 6 {
		_reduceAfterSub(t, false)
	} else {
		_reduceAfterSub(t, true)
	}

	r = x
	movq("res+0(FP)", r)

	_mov(t, r, 0, nbWords)

	ret()

}
