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

import (
	"fmt"

	"github.com/consensys/bavard/amd64"
)

func (f *FFAmd64) generateReduce() {
	stackSize := 0
	if f.NbWords > SmallModulus {
		stackSize = f.NbWords * 8
	}
	registers := f.FnHeader("reduce", stackSize, 8)

	// registers
	r := registers.Pop()
	t := registers.PopN(f.NbWords)

	f.MOVQ("res+0(FP)", r)

	f.Mov(r, t)
	f.Reduce(&registers, t, r)
	f.RET()
}

func (f *FFAmd64) Reduce(registers *amd64.Registers, t []amd64.Register, result interface{}, rOffset ...int) {
	if f.NbWords > SmallModulus {
		f.reduceLarge(t, result, rOffset...)
		return
	}
	offset := 0
	if len(rOffset) > 0 {
		offset = rOffset[0]
	}

	// u = t - q
	u := registers.PopN(f.NbWords)

	if rm, ok := result.(amd64.Register); ok && offset == 0 {
		f.reduceElement(rm, t, u)
	} else {
		f.copyElement(t, u)
		for i := 0; i < f.NbWords; i++ {
			if i == 0 {
				f.SUBQ(f.qAt(i), u[i])
			} else {
				f.SBBQ(f.qAt(i), u[i])
			}
		}

		// conditional move of u into t (if we have a borrow we need to return t - q)
		for i := 0; i < f.NbWords; i++ {
			f.CMOVQCC(u[i], t[i])
		}

		// return t
		f.Mov(t, result, 0, offset)
	}

	registers.Push(u...)
}

func (f *FFAmd64) reduceLarge(t []amd64.Register, result interface{}, rOffset ...int) {
	// u = t - q
	u := make([]string, f.NbWords)

	for i := 0; i < f.NbWords; i++ {
		// use stack
		u[i] = fmt.Sprintf("t%d-%d(SP)", i, 8+i*8)
		f.MOVQ(t[i], u[i])

		if i == 0 {
			f.SUBQ(f.qAt(i), t[i])
		} else {
			f.SBBQ(f.qAt(i), t[i])
		}
	}

	// conditional move of u into t (if we have a borrow we need to return t - q)
	for i := 0; i < f.NbWords; i++ {
		f.CMOVQCS(u[i], t[i])
	}

	offset := 0
	if len(rOffset) > 0 {
		offset = rOffset[0]
	}
	// return t
	f.Mov(t, result, 0, offset)
}
