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

// field modulus q =
//
// 21888242871839275222246405745257275088696311157297823662689037894645226208583
// Code generated by goff DO NOT EDIT
// Element are assumed to be in Montgomery form in all methods

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

// Package bn256 (generated by goff) contains field arithmetics operations
package bn256

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"math/big"
	"math/bits"
	"sync"

	"unsafe"
)

// Element represents a field element stored on 4 words (uint64)
// Element are assumed to be in Montgomery form in all methods
type Element [4]uint64

// ElementLimbs number of 64 bits words needed to represent Element
const ElementLimbs = 4

// ElementBits number bits needed to represent Element
const ElementBits = 254

// SetUint64 z = v, sets z LSB to v (non-Montgomery form) and convert z to Montgomery form
func (z *Element) SetUint64(v uint64) *Element {
	z[0] = v
	z[1] = 0
	z[2] = 0
	z[3] = 0
	return z.ToMont()
}

// Set z = x
func (z *Element) Set(x *Element) *Element {
	z[0] = x[0]
	z[1] = x[1]
	z[2] = x[2]
	z[3] = x[3]
	return z
}

// SetZero z = 0
func (z *Element) SetZero() *Element {
	z[0] = 0
	z[1] = 0
	z[2] = 0
	z[3] = 0
	return z
}

// SetOne z = 1 (in Montgomery form)
func (z *Element) SetOne() *Element {
	z[0] = 15230403791020821917
	z[1] = 754611498739239741
	z[2] = 7381016538464732716
	z[3] = 1011752739694698287
	return z
}

// Neg z = q - x
func (z *Element) Neg(x *Element) *Element {
	if x.IsZero() {
		return z.SetZero()
	}
	var borrow uint64
	z[0], borrow = bits.Sub64(4332616871279656263, x[0], 0)
	z[1], borrow = bits.Sub64(10917124144477883021, x[1], borrow)
	z[2], borrow = bits.Sub64(13281191951274694749, x[2], borrow)
	z[3], _ = bits.Sub64(3486998266802970665, x[3], borrow)
	return z
}

// Div z = x*y^-1 mod q
func (z *Element) Div(x, y *Element) *Element {
	var yInv Element
	yInv.Inverse(y)
	z.Mul(x, &yInv)
	return z
}

// Equal returns z == x
func (z *Element) Equal(x *Element) bool {
	return (z[3] == x[3]) && (z[2] == x[2]) && (z[1] == x[1]) && (z[0] == x[0])
}

// IsZero returns z == 0
func (z *Element) IsZero() bool {
	return (z[3] | z[2] | z[1] | z[0]) == 0
}

// field modulus stored as big.Int
var _elementModulusBigInt big.Int
var onceelementModulus sync.Once

func elementModulusBigInt() *big.Int {
	onceelementModulus.Do(func() {
		_elementModulusBigInt.SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	})
	return &_elementModulusBigInt
}

// Inverse z = x^-1 mod q
// Algorithm 16 in "Efficient Software-Implementation of Finite Fields with Applications to Cryptography"
// if x == 0, sets and returns z = x
func (z *Element) Inverse(x *Element) *Element {
	if x.IsZero() {
		return z.Set(x)
	}

	// initialize u = q
	var u = Element{
		4332616871279656263,
		10917124144477883021,
		13281191951274694749,
		3486998266802970665,
	}

	// initialize s = r^2
	var s = Element{
		17522657719365597833,
		13107472804851548667,
		5164255478447964150,
		493319470278259999,
	}

	// r = 0
	r := Element{}

	v := *x

	var carry, borrow, t, t2 uint64
	var bigger, uIsOne, vIsOne bool

	for !uIsOne && !vIsOne {
		for v[0]&1 == 0 {

			// v = v >> 1
			t2 = v[3] << 63
			v[3] >>= 1
			t = t2
			t2 = v[2] << 63
			v[2] = (v[2] >> 1) | t
			t = t2
			t2 = v[1] << 63
			v[1] = (v[1] >> 1) | t
			t = t2
			v[0] = (v[0] >> 1) | t

			if s[0]&1 == 1 {

				// s = s + q
				s[0], carry = bits.Add64(s[0], 4332616871279656263, 0)
				s[1], carry = bits.Add64(s[1], 10917124144477883021, carry)
				s[2], carry = bits.Add64(s[2], 13281191951274694749, carry)
				s[3], _ = bits.Add64(s[3], 3486998266802970665, carry)

			}

			// s = s >> 1
			t2 = s[3] << 63
			s[3] >>= 1
			t = t2
			t2 = s[2] << 63
			s[2] = (s[2] >> 1) | t
			t = t2
			t2 = s[1] << 63
			s[1] = (s[1] >> 1) | t
			t = t2
			s[0] = (s[0] >> 1) | t

		}
		for u[0]&1 == 0 {

			// u = u >> 1
			t2 = u[3] << 63
			u[3] >>= 1
			t = t2
			t2 = u[2] << 63
			u[2] = (u[2] >> 1) | t
			t = t2
			t2 = u[1] << 63
			u[1] = (u[1] >> 1) | t
			t = t2
			u[0] = (u[0] >> 1) | t

			if r[0]&1 == 1 {

				// r = r + q
				r[0], carry = bits.Add64(r[0], 4332616871279656263, 0)
				r[1], carry = bits.Add64(r[1], 10917124144477883021, carry)
				r[2], carry = bits.Add64(r[2], 13281191951274694749, carry)
				r[3], _ = bits.Add64(r[3], 3486998266802970665, carry)

			}

			// r = r >> 1
			t2 = r[3] << 63
			r[3] >>= 1
			t = t2
			t2 = r[2] << 63
			r[2] = (r[2] >> 1) | t
			t = t2
			t2 = r[1] << 63
			r[1] = (r[1] >> 1) | t
			t = t2
			r[0] = (r[0] >> 1) | t

		}

		// v >= u
		bigger = !(v[3] < u[3] || (v[3] == u[3] && (v[2] < u[2] || (v[2] == u[2] && (v[1] < u[1] || (v[1] == u[1] && (v[0] < u[0])))))))

		if bigger {

			// v = v - u
			v[0], borrow = bits.Sub64(v[0], u[0], 0)
			v[1], borrow = bits.Sub64(v[1], u[1], borrow)
			v[2], borrow = bits.Sub64(v[2], u[2], borrow)
			v[3], _ = bits.Sub64(v[3], u[3], borrow)

			// r >= s
			bigger = !(r[3] < s[3] || (r[3] == s[3] && (r[2] < s[2] || (r[2] == s[2] && (r[1] < s[1] || (r[1] == s[1] && (r[0] < s[0])))))))

			if bigger {

				// s = s + q
				s[0], carry = bits.Add64(s[0], 4332616871279656263, 0)
				s[1], carry = bits.Add64(s[1], 10917124144477883021, carry)
				s[2], carry = bits.Add64(s[2], 13281191951274694749, carry)
				s[3], _ = bits.Add64(s[3], 3486998266802970665, carry)

			}

			// s = s - r
			s[0], borrow = bits.Sub64(s[0], r[0], 0)
			s[1], borrow = bits.Sub64(s[1], r[1], borrow)
			s[2], borrow = bits.Sub64(s[2], r[2], borrow)
			s[3], _ = bits.Sub64(s[3], r[3], borrow)

		} else {

			// u = u - v
			u[0], borrow = bits.Sub64(u[0], v[0], 0)
			u[1], borrow = bits.Sub64(u[1], v[1], borrow)
			u[2], borrow = bits.Sub64(u[2], v[2], borrow)
			u[3], _ = bits.Sub64(u[3], v[3], borrow)

			// s >= r
			bigger = !(s[3] < r[3] || (s[3] == r[3] && (s[2] < r[2] || (s[2] == r[2] && (s[1] < r[1] || (s[1] == r[1] && (s[0] < r[0])))))))

			if bigger {

				// r = r + q
				r[0], carry = bits.Add64(r[0], 4332616871279656263, 0)
				r[1], carry = bits.Add64(r[1], 10917124144477883021, carry)
				r[2], carry = bits.Add64(r[2], 13281191951274694749, carry)
				r[3], _ = bits.Add64(r[3], 3486998266802970665, carry)

			}

			// r = r - s
			r[0], borrow = bits.Sub64(r[0], s[0], 0)
			r[1], borrow = bits.Sub64(r[1], s[1], borrow)
			r[2], borrow = bits.Sub64(r[2], s[2], borrow)
			r[3], _ = bits.Sub64(r[3], s[3], borrow)

		}
		uIsOne = (u[0] == 1) && (u[3]|u[2]|u[1]) == 0
		vIsOne = (v[0] == 1) && (v[3]|v[2]|v[1]) == 0
	}

	if uIsOne {
		z.Set(&r)
	} else {
		z.Set(&s)
	}

	return z
}

// SetRandom sets z to a random element < q
func (z *Element) SetRandom() *Element {
	bytes := make([]byte, 32)
	io.ReadFull(rand.Reader, bytes)
	z[0] = binary.BigEndian.Uint64(bytes[0:8])
	z[1] = binary.BigEndian.Uint64(bytes[8:16])
	z[2] = binary.BigEndian.Uint64(bytes[16:24])
	z[3] = binary.BigEndian.Uint64(bytes[24:32])
	z[3] %= 3486998266802970665

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}

	return z
}

// Add z = x + y mod q
func (z *Element) Add(x, y *Element) *Element {
	var carry uint64

	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], _ = bits.Add64(x[3], y[3], carry)

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z
}

// AddAssign z = z + x mod q
func (z *Element) AddAssign(x *Element) *Element {
	var carry uint64

	z[0], carry = bits.Add64(z[0], x[0], 0)
	z[1], carry = bits.Add64(z[1], x[1], carry)
	z[2], carry = bits.Add64(z[2], x[2], carry)
	z[3], _ = bits.Add64(z[3], x[3], carry)

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z
}

// Double z = x + x mod q, aka Lsh 1
func (z *Element) Double(x *Element) *Element {
	var carry uint64

	z[0], carry = bits.Add64(x[0], x[0], 0)
	z[1], carry = bits.Add64(x[1], x[1], carry)
	z[2], carry = bits.Add64(x[2], x[2], carry)
	z[3], _ = bits.Add64(x[3], x[3], carry)

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z
}

// Sub  z = x - y mod q
func (z *Element) Sub(x, y *Element) *Element {
	var b uint64
	z[0], b = bits.Sub64(x[0], y[0], 0)
	z[1], b = bits.Sub64(x[1], y[1], b)
	z[2], b = bits.Sub64(x[2], y[2], b)
	z[3], b = bits.Sub64(x[3], y[3], b)
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], 4332616871279656263, 0)
		z[1], c = bits.Add64(z[1], 10917124144477883021, c)
		z[2], c = bits.Add64(z[2], 13281191951274694749, c)
		z[3], _ = bits.Add64(z[3], 3486998266802970665, c)
	}
	return z
}

// SubAssign  z = z - x mod q
func (z *Element) SubAssign(x *Element) *Element {
	var b uint64
	z[0], b = bits.Sub64(z[0], x[0], 0)
	z[1], b = bits.Sub64(z[1], x[1], b)
	z[2], b = bits.Sub64(z[2], x[2], b)
	z[3], b = bits.Sub64(z[3], x[3], b)
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], 4332616871279656263, 0)
		z[1], c = bits.Add64(z[1], 10917124144477883021, c)
		z[2], c = bits.Add64(z[2], 13281191951274694749, c)
		z[3], _ = bits.Add64(z[3], 3486998266802970665, c)
	}
	return z
}

// Exp z = x^exponent mod q
// (not optimized)
// exponent (non-montgomery form) is ordered from least significant word to most significant word
func (z *Element) Exp(x Element, exponent ...uint64) *Element {
	r := 0
	msb := 0
	for i := len(exponent) - 1; i >= 0; i-- {
		if exponent[i] == 0 {
			r++
		} else {
			msb = (i * 64) + bits.Len64(exponent[i])
			break
		}
	}
	exponent = exponent[:len(exponent)-r]
	if len(exponent) == 0 {
		return z.SetOne()
	}

	z.Set(&x)

	l := msb - 2
	for i := l; i >= 0; i-- {
		z.Square(z)
		if exponent[i/64]&(1<<uint(i%64)) != 0 {
			z.MulAssign(&x)
		}
	}
	return z
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *Element) FromMont() *Element {

	// the following lines implement z = z * 1
	// with a modified CIOS montgomery multiplication
	{
		// m = z[0]n'[0] mod W
		m := z[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, z[0])
		C, z[0] = madd2(m, 10917124144477883021, z[1], C)
		C, z[1] = madd2(m, 13281191951274694749, z[2], C)
		C, z[2] = madd2(m, 3486998266802970665, z[3], C)
		z[3] = C
	}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, z[0])
		C, z[0] = madd2(m, 10917124144477883021, z[1], C)
		C, z[1] = madd2(m, 13281191951274694749, z[2], C)
		C, z[2] = madd2(m, 3486998266802970665, z[3], C)
		z[3] = C
	}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, z[0])
		C, z[0] = madd2(m, 10917124144477883021, z[1], C)
		C, z[1] = madd2(m, 13281191951274694749, z[2], C)
		C, z[2] = madd2(m, 3486998266802970665, z[3], C)
		z[3] = C
	}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, z[0])
		C, z[0] = madd2(m, 10917124144477883021, z[1], C)
		C, z[1] = madd2(m, 13281191951274694749, z[2], C)
		C, z[2] = madd2(m, 3486998266802970665, z[3], C)
		z[3] = C
	}

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z
}

// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
func (z *Element) ToMont() *Element {
	var rSquare = Element{
		17522657719365597833,
		13107472804851548667,
		5164255478447964150,
		493319470278259999,
	}
	return z.MulAssign(&rSquare)
}

// ToRegular returns z in regular form (doesn't mutate z)
func (z Element) ToRegular() Element {
	return *z.FromMont()
}

// String returns the string form of an Element in Montgomery form
func (z *Element) String() string {
	var _z big.Int
	return z.ToBigIntRegular(&_z).String()
}

// ToBigInt returns z as a big.Int in Montgomery form
func (z *Element) ToBigInt(res *big.Int) *big.Int {
	bits := (*[4]big.Word)(unsafe.Pointer(z))
	return res.SetBits(bits[:])
}

// ToBigIntRegular returns z as a big.Int in regular form
func (z Element) ToBigIntRegular(res *big.Int) *big.Int {
	z.FromMont()
	bits := (*[4]big.Word)(unsafe.Pointer(&z))
	return res.SetBits(bits[:])
}

// SetBigInt sets z to v (regular form) and returns z in Montgomery form
func (z *Element) SetBigInt(v *big.Int) *Element {
	z.SetZero()

	zero := big.NewInt(0)
	q := elementModulusBigInt()

	// copy input
	vv := new(big.Int).Set(v)

	// while v < 0, v+=q
	for vv.Cmp(zero) == -1 {
		vv.Add(vv, q)
	}
	// while v > q, v-=q
	for vv.Cmp(q) == 1 {
		vv.Sub(vv, q)
	}
	// if v == q, return 0
	if vv.Cmp(q) == 0 {
		return z
	}
	// v should
	vBits := vv.Bits()
	for i := 0; i < len(vBits); i++ {
		z[i] = uint64(vBits[i])
	}
	return z.ToMont()
}

// SetString creates a big.Int with s (in base 10) and calls SetBigInt on z
func (z *Element) SetString(s string) *Element {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("Element.SetString failed -> can't parse number in base10 into a big.Int")
	}
	return z.SetBigInt(x)
}

// Mul z = x * y mod q
func (z *Element) Mul(x, y *Element) *Element {

	var t [4]uint64
	var c [3]uint64
	{
		// round 0
		v := x[0]
		c[1], c[0] = bits.Mul64(v, y[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd1(v, y[1], c[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd1(v, y[2], c[1])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd1(v, y[3], c[1])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 1
		v := x[1]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 2
		v := x[2]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 3
		v := x[3]
		c[1], c[0] = madd1(v, y[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, y[1], c[1], t[1])
		c[2], z[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, y[2], c[1], t[2])
		c[2], z[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, y[3], c[1], t[3])
		z[3], z[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z
}

// MulAssign z = z * x mod q
func (z *Element) MulAssign(x *Element) *Element {

	var t [4]uint64
	var c [3]uint64
	{
		// round 0
		v := z[0]
		c[1], c[0] = bits.Mul64(v, x[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd1(v, x[1], c[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd1(v, x[2], c[1])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd1(v, x[3], c[1])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 1
		v := z[1]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 2
		v := z[2]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], t[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], t[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		t[3], t[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}
	{
		// round 3
		v := z[3]
		c[1], c[0] = madd1(v, x[0], t[0])
		m := c[0] * 9786893198990664585
		c[2] = madd0(m, 4332616871279656263, c[0])
		c[1], c[0] = madd2(v, x[1], c[1], t[1])
		c[2], z[0] = madd2(m, 10917124144477883021, c[2], c[0])
		c[1], c[0] = madd2(v, x[2], c[1], t[2])
		c[2], z[1] = madd2(m, 13281191951274694749, c[2], c[0])
		c[1], c[0] = madd2(v, x[3], c[1], t[3])
		z[3], z[2] = madd3(m, 3486998266802970665, c[0], c[2], c[1])
	}

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z
}

func (z *Element) Legendre() int {
	var l Element
	// z^((p-1)/2)
	l.Exp(*z,
		11389680472494603939,
		14681934109093717318,
		15863968012492123182,
		1743499133401485332,
	)

	if l.IsZero() {
		return 0
	}

	// if l == 1
	if (l[3] == 1011752739694698287) && (l[2] == 7381016538464732716) && (l[1] == 754611498739239741) && (l[0] == 15230403791020821917) {
		return 1
	}
	return -1
}

// Sqrt z = √x mod q
// if the square root doesn't exist (x is not a square mod q)
// Sqrt leaves z unchanged and returns nil
func (z *Element) Sqrt(x *Element) *Element {
	switch x.Legendre() {
	case -1:
		return nil
	case 0:
		return z.SetZero()
	case 1:
		break
	}
	// q ≡ 3 (mod 4)
	// using  z ≡ ± x^((p+1)/4) (mod q)
	return z.Exp(*x,
		5694840236247301970,
		7340967054546858659,
		7931984006246061591,
		871749566700742666,
	)
}

// Square z = x * x mod q
func (z *Element) Square(x *Element) *Element {

	var p [4]uint64

	var u, v uint64
	{
		// round 0
		u, p[0] = bits.Mul64(x[0], x[0])
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		var t uint64
		t, u, v = madd1sb(x[0], x[1], u)
		C, p[0] = madd2(m, 10917124144477883021, v, C)
		t, u, v = madd1s(x[0], x[2], t, u)
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd1s(x[0], x[3], t, u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 1
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		u, v = madd1(x[1], x[1], p[1])
		C, p[0] = madd2(m, 10917124144477883021, v, C)
		var t uint64
		t, u, v = madd2sb(x[1], x[2], p[2], u)
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd2s(x[1], x[3], p[3], t, u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 2
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		C, p[0] = madd2(m, 10917124144477883021, p[1], C)
		u, v = madd1(x[2], x[2], p[2])
		C, p[1] = madd2(m, 13281191951274694749, v, C)
		_, u, v = madd2sb(x[2], x[3], p[3], u)
		p[3], p[2] = madd3(m, 3486998266802970665, v, C, u)
	}
	{
		// round 3
		m := p[0] * 9786893198990664585
		C := madd0(m, 4332616871279656263, p[0])
		C, z[0] = madd2(m, 10917124144477883021, p[1], C)
		C, z[1] = madd2(m, 13281191951274694749, p[2], C)
		u, v = madd1(x[3], x[3], p[3])
		z[3], z[2] = madd3(m, 3486998266802970665, v, C, u)
	}

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 3486998266802970665 || (z[3] == 3486998266802970665 && (z[2] < 13281191951274694749 || (z[2] == 13281191951274694749 && (z[1] < 10917124144477883021 || (z[1] == 10917124144477883021 && (z[0] < 4332616871279656263))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 4332616871279656263, 0)
		z[1], b = bits.Sub64(z[1], 10917124144477883021, b)
		z[2], b = bits.Sub64(z[2], 13281191951274694749, b)
		z[3], _ = bits.Sub64(z[3], 3486998266802970665, b)
	}
	return z
}
