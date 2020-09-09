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

import "fmt"

const dx = register("DX")
const ax = register("AX")

type lbl string
type register string

func (r *register) at(wordOffset int) string {
	return fmt.Sprintf("%d(%s)", wordOffset*8, string(*r))
}

func availableRegisters() int {
	return len(builder.registers)
}

func popRegister() register {
	r := builder.registers[0]
	builder.registers = builder.registers[1:]
	return r
}

func popRegisters(n int) []register {
	toReturn := make([]register, n)
	for i := 0; i < n; i++ {
		toReturn[i] = popRegister()
	}
	return toReturn
}

func pushRegister(r ...register) {
	builder.registers = append(builder.registers, r...)
}

var staticRegisters = []register{
	"AX",
	"DX",
	"CX",
	"BX",
	"BP",
	"SI",
	"DI",
	"R8",
	"R9",
	"R10",
	"R11",
	"R12",
	"R13",
	"R14",
	"R15",
}

var labelCounter = 0

func newLabel() lbl {
	labelCounter++
	return lbl(fmt.Sprintf("l%d", labelCounter))
}
