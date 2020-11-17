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

// Package amd64 contains syntactic sugar to generate amd64 assembly code in goff
package amd64

import (
	"fmt"
	"io"

	"github.com/consensys/bavard"

	. "github.com/consensys/bavard/amd64"
	"github.com/consensys/goff/field"
)

const SmallModulus = 6

func newFFAmd64(w io.Writer, F *field.Field) *ffAmd64 {
	Lock(w)
	return &ffAmd64{F}
}

type ffAmd64 struct {
	*field.Field
}

func (f *ffAmd64) qAt(index int) string {
	return fmt.Sprintf("·q%s+%d(SB)", f.ElementName, index*8)
}

func (f *ffAmd64) qInv0() string {
	return fmt.Sprintf("·q%sInv0(SB)", f.ElementName)
}

// Generate generates assembly code for the base field provided to goff
// see internal/templates/ops*
func Generate(w io.Writer, F *field.Field) error {
	f := newFFAmd64(w, F)
	defer Unlock()
	WriteLn(bavard.Apache2Header("ConsenSys Software Inc.", 2020))

	WriteLn("#include \"textflag.h\"")
	WriteLn("#include \"funcdata.h\"")

	// mul
	f.generateMul()
	// square
	f.generateSquare()

	// from mont
	f.generateFromMont()

	// reduce
	f.generateReduce()

	// add
	f.generateAdd()

	// sub
	f.generateSub()

	// double
	f.generateDouble()

	// neg
	f.generateNeg()

	return nil
}
