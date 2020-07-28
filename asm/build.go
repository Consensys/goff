package asm

import (
	"fmt"
	"io"
)

const smallModulus = 6

var (
	nbWords          int
	nbWordsLastIndex int
	elementName      string
	modulus          []uint64
	noCarrySquare    bool
	builder          *assembly
)

type assembly struct {
	writer    io.Writer
	registers []register
}

// NewBuilder returns a builder object to help generated assembly code for some operations
func NewBuilder(w io.Writer, _elementName string, _nbWords int, _q []uint64, _noCarrySquare bool) *assembly {
	b := &assembly{
		writer:    w,
		registers: make([]register, len(staticRegisters)),
	}
	elementName = _elementName
	modulus = _q
	nbWords = _nbWords
	nbWordsLastIndex = _nbWords - 1
	noCarrySquare = _noCarrySquare
	copy(b.registers, staticRegisters)
	builder = b
	return b
}

func qAt(index int) string {
	return fmt.Sprintf("·q%s+%d(SB)", elementName, index*8)
}

func qInv0() string {
	return fmt.Sprintf("·q%sInv0(SB)", elementName)
}

// Build ...
func (asm *assembly) Build() error {

	apache2 := `
// Copyright %d %s
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
	`
	license := fmt.Sprintf(apache2, 2020, "ConsenSys Software Inc.")
	writeLn(license)

	writeLn("#include \"textflag.h\"")
	writeLn("#include \"funcdata.h\"")

	// mul
	generateMul()
	// square
	generateSquare()

	// from mont
	generateFromMont()

	// reduce
	generateReduce()

	// add
	generateAdd()
	generateAddE2()

	// sub
	generateSub()
	generateSubE2()

	// double
	generateDouble()
	generateDoubleE2()

	// neg
	generateNeg()
	generateNegE2()

	return nil
}
