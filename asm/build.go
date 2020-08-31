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

// GenerateAssembly generates assembly code for the base field provided to goff
// see internal/templates/ops*
func (asm *assembly) GenerateAssembly() error {

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

	// sub
	generateSub()

	// double
	generateDouble()

	// neg
	generateNeg()

	return nil
}

// SpecialCurve is likely temporary --> this should move into gurvy package in the future
type SpecialCurve int

// SpecialCurve enum to generate curve specific tower of extension functions.
// note that this is likely going to move in gurvy at some point.
const (
	NONE SpecialCurve = iota
	BN256
	BLS381
)

// GenerateTowerAssembly will generate assembly function for the tower of extension provided to goff
// note that this is likely going to move in gurvy at some point
func (asm *assembly) GenerateTowerAssembly(specialCurve SpecialCurve) error {

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

	generateAddE2()
	generateSubE2()
	generateDoubleE2()
	generateNegE2()

	switch specialCurve {
	case BN256:
		generateSquareE2BN256()
		generateMulE2BN256()
		generateMulByNonResidueE2BN256()
	case BLS381:
		fmt.Println("bls381 e2 square, mul and mul by non residue not implemented in asm")
	}

	return nil
}
