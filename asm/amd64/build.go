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

	"github.com/consensys/bavard/amd64"
	"github.com/consensys/goff/field"
)

const SmallModulus = 6

func NewFFAmd64(w io.Writer, F *field.Field) *FFAmd64 {
	return &FFAmd64{F, amd64.NewAmd64(w)}
}

type FFAmd64 struct {
	*field.Field
	*amd64.Amd64
}

func (f *FFAmd64) qAt(index int) string {
	return fmt.Sprintf("·q%s+%d(SB)", f.ElementName, index*8)
}

func (f *FFAmd64) qInv0() string {
	return fmt.Sprintf("·q%sInv0(SB)", f.ElementName)
}

// Generate generates assembly code for the base field provided to goff
// see internal/templates/ops*
func Generate(w io.Writer, F *field.Field) error {
	f := NewFFAmd64(w, F)
	f.WriteLn(bavard.Apache2Header("ConsenSys Software Inc.", 2020))

	f.WriteLn("#include \"textflag.h\"")
	f.WriteLn("#include \"funcdata.h\"")

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
