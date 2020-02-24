// Copyright 2019 ConsenSys AG
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

package cmd

import (
	"math/big"
)

type field struct {
	PackageName          string
	ElementName          string
	Modulus              string
	NbWords              int
	NbBits               int
	NbWordsLastIndex     int
	NbWordsIndexesNoZero []int
	NbWordsIndexesFull   []int
	IdxFIPS              []int
	Q                    []uint64
	QInverse             []uint64
	RSquare              []uint64
	One                  []uint64
	NoCarry              bool
	NoCarrySquare        bool // used if NoCarry is set, but some op may overflow in square optimization
	Benches              bool
	Version              string
}

// -------------------------------------------------------------------------------------------------
// Field data precompute functions
func newField(packageName, elementName, modulus string, benches bool) (*field, error) {
	// parse modulus
	var bModulus big.Int
	if _, ok := bModulus.SetString(modulus, 10); !ok {
		return nil, errParseModulus
	}

	// field info
	F := &field{
		PackageName: packageName,
		ElementName: elementName,
		Modulus:     modulus,
		Benches:     benches,
		Version:     buildString(),
	}

	// pre compute field constants
	F.NbBits = bModulus.BitLen()
	F.NbWords = len(bModulus.Bits())
	if F.NbWords < 2 {
		return nil, errUnsupportedModulus
	}

	F.NbWordsLastIndex = F.NbWords - 1
	F.Q = make([]uint64, F.NbWords)
	F.QInverse = make([]uint64, F.NbWords)
	F.RSquare = make([]uint64, F.NbWords)
	F.One = make([]uint64, F.NbWords)

	// set q from big int repr
	for i, v := range bModulus.Bits() {
		F.Q[i] = (uint64)(v)
	}

	//  setting qInverse
	_r := big.NewInt(1)
	_r.Lsh(_r, uint(F.NbWords)*64)
	_rInv := big.NewInt(1)
	_qInv := big.NewInt(0)
	extendedEuclideanAlgo(_r, &bModulus, _rInv, _qInv)
	_qInv.Mod(_qInv, _r)
	for i, v := range _qInv.Bits() {
		F.QInverse[i] = (uint64)(v)
	}

	//  rsquare
	_rSquare := big.NewInt(2)
	exponent := big.NewInt(int64(F.NbWords) * 64 * 2)
	_rSquare.Exp(_rSquare, exponent, &bModulus)
	_rSquareBits := _rSquare.Bits()
	for i := 0; i < len(_rSquareBits); i++ {
		F.RSquare[i] = uint64(_rSquareBits[i])
	}

	var one big.Int
	one.SetUint64(1)
	one.Lsh(&one, uint(F.NbWords)*64).Mod(&one, &bModulus)
	_oneBits := one.Bits()
	for i := 0; i < len(_oneBits); i++ {
		F.One[i] = uint64(_oneBits[i])
	}

	// indexes (template helpers)
	F.NbWordsIndexesFull = make([]int, F.NbWords)
	F.NbWordsIndexesNoZero = make([]int, F.NbWords-1)
	for i := 0; i < F.NbWords; i++ {
		F.NbWordsIndexesFull[i] = i
		if i > 0 {
			F.NbWordsIndexesNoZero[i-1] = i
		}
	}

	// See https:// TODO blog post link
	// if the last word of the modulus is smaller or equal to B,
	// we can simplify the montgomery multiplication
	const B = (^uint64(0) >> 1) - 1
	F.NoCarry = (F.Q[len(F.Q)-1] <= B) && F.NbWords <= 12
	const BSquare = (^uint64(0) >> 2)
	F.NoCarrySquare = F.Q[len(F.Q)-1] <= BSquare

	for i := F.NbWords; i <= 2*F.NbWords-2; i++ {
		F.IdxFIPS = append(F.IdxFIPS, i)
	}

	return F, nil
}

// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
// r > q, modifies rinv and qinv such that rinv.r - qinv.q = 1
func extendedEuclideanAlgo(r, q, rInv, qInv *big.Int) {
	var s1, s2, t1, t2, qi, tmpMuls, riPlusOne, tmpMult, a, b big.Int
	t1.SetUint64(1)
	rInv.Set(big.NewInt(1))
	qInv.Set(big.NewInt(0))
	a.Set(r)
	b.Set(q)

	// r_i+1 = r_i-1 - q_i.r_i
	// s_i+1 = s_i-1 - q_i.s_i
	// t_i+1 = t_i-1 - q_i.s_i
	for b.Sign() > 0 {
		qi.Div(&a, &b)
		riPlusOne.Mod(&a, &b)

		tmpMuls.Mul(&s1, &qi)
		tmpMult.Mul(&t1, &qi)

		s2.Set(&s1)
		t2.Set(&t1)

		s1.Sub(rInv, &tmpMuls)
		t1.Sub(qInv, &tmpMult)
		rInv.Set(&s2)
		qInv.Set(&t2)

		a.Set(&b)
		b.Set(&riPlusOne)
	}
	qInv.Neg(qInv)
}
