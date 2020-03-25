package element

const Arith = `
import (
	"math/bits"
)

func madd(a, b, t,u,v uint64) (uint64, uint64, uint64) {
	var carry uint64
	hi, lo := bits.Mul64(a, b)
	v, carry = bits.Add64(lo, v, 0)
	u, carry = bits.Add64(hi, u, carry)
	t, _ = bits.Add64(t, 0, carry)
	return t, u, v 
}

// madd0 hi = a*b + c (discards lo bits)
func madd0(a, b, c uint64) (hi uint64) {
	var carry, lo uint64
	hi, lo = bits.Mul64(a, b)
	_, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

// madd1 hi, lo = a*b + c
func madd1(a, b, c uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

// madd2 hi, lo = a*b + c + d
func madd2(a, b, c, d uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	c, carry = bits.Add64(c, d, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

// madd2s superhi, hi, lo = 2*a*b + c + d + e
func madd2s(a, b, c, d, e uint64) (superhi, hi, lo uint64) {
	var carry, sum uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)

	sum, carry = bits.Add64(c, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, sum, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	hi, _ = bits.Add64(hi, 0, d)
	return
}

func madd1s(a, b,  d, e uint64) (superhi, hi, lo uint64) {
	var carry uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)
	lo, carry = bits.Add64(lo, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	hi, _ = bits.Add64(hi, 0, d)
	return
}


func madd2sb(a, b, c, e uint64) (superhi, hi, lo uint64) {
	var carry, sum uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)

	sum, carry = bits.Add64(c, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, sum, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

func madd1sb(a, b, e uint64) (superhi, hi, lo uint64) {
	var carry uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)
	lo, carry = bits.Add64(lo, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}



func madd3(a, b, c, d, e uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	c, carry = bits.Add64(c, d, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, e, carry)
	return
}


`
