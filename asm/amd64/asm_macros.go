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

import "github.com/consensys/bavard/amd64"

func (f *FFAmd64) Mov(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {
	case []uint64:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1[i+o1], c2[i+o2])
			}
		}
	case amd64.Register:
		switch c2 := i2.(type) {
		case amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1.At(i+o1), c2.At(i+o2))
			}
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1.At(i+o1), c2[i+o2])
			}
		default:
			panic("unsupported")
		}
	case []amd64.Register:
		switch c2 := i2.(type) {
		case amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1[i+o1], c2.At(i+o2))
			}
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1[i+o1], c2[i+o2])
			}
		default:
			panic("unsupported")
		}
	default:
		panic("unsupported")
	}

}

func (f *FFAmd64) Add(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {

	case amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.ADDQ(c1.At(i+o1), c2[i+o2])
				} else {
					f.ADCQ(c1.At(i+o1), c2[i+o2])
				}
			}
		}
	case []amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.ADDQ(c1[i+o1], c2[i+o2])
				} else {
					f.ADCQ(c1[i+o1], c2[i+o2])
				}
			}
		}
	default:
		panic("unsupported")
	}
}

func (f *FFAmd64) Sub(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {

	case amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.SUBQ(c1.At(i+o1), c2[i+o2])
				} else {
					f.SBBQ(c1.At(i+o1), c2[i+o2])
				}
			}
		}
	case []amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.SUBQ(c1[i+o1], c2[i+o2])
				} else {
					f.SBBQ(c1[i+o1], c2[i+o2])
				}
			}
		}
	default:
		panic("unsupported")
	}
}
