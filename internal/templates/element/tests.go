package element

const Test = `

import (
    "crypto/rand"
	"math/big"
	"math/bits"
	"testing"
	"fmt"
    mrand "math/rand"
)

func Test{{toUpper .ElementName}}CorrectnessAgainstBigInt(t *testing.T) {
    modulus, _ := new(big.Int).SetString("{{.Modulus}}", 10)
	cmpEandB := func(e *{{.ElementName}}, b *big.Int, name string) {
		var _e big.Int
		if e.FromMont().ToBigInt(&_e).Cmp(b) != 0 {
			t.Fatal(name, "failed")
		}
	}
	var modulusMinusOne, one big.Int
	one.SetUint64(1)

	modulusMinusOne.Sub(modulus, &one)

	var n int
	if testing.Short() {
		n = 10
	} else {
		n = 500
	}

    for i := 0; i < n; i++ {

        // sample 2 random big int
        b1, _ := rand.Int(rand.Reader, modulus)
        b2, _ := rand.Int(rand.Reader, modulus)
        rExp := mrand.Uint64()
        

        // adding edge cases
        // TODO need more edge cases
        switch i {
        case 0:
            rExp = 0
            b1.SetUint64(0)
        case 1:
            b2.SetUint64(0)
        case 2:
            b1.SetUint64(0)
            b2.SetUint64(0)
        case 3:
            rExp = 0
        case 4:
            rExp = 1
        case 5:
			rExp = ^uint64(0) // max uint
		case 6:
			rExp = 2
			b1.Set(&modulusMinusOne)
		case 7:
			b2.Set(&modulusMinusOne)
		case 8:
			b1.Set(&modulusMinusOne)
			b2.Set(&modulusMinusOne)
        }

        rbExp := new(big.Int).SetUint64(rExp)

        var bMul, bAdd, bSub, bDiv, bNeg, bLsh, bInv, bExp, bExp2,  bSquare big.Int

        // e1 = mont(b1), e2 = mont(b2)
        var e1, e2, eMul,  eAdd, eSub, eDiv, eNeg, eLsh, eInv, eExp, eExp2, eSquare, eMulAssign, eSubAssign, eAddAssign {{.ElementName}}
        e1.SetBigInt(b1)
        e2.SetBigInt(b2)

        // (e1*e2).FromMont() === b1*b2 mod q ... etc
        eSquare.Square(&e1)
		eMul.Mul(&e1, &e2)
        eMulAssign.Set(&e1)
        eMulAssign.MulAssign(&e2)
        eAdd.Add(&e1, &e2)
        eAddAssign.Set(&e1)
        eAddAssign.AddAssign(&e2)
        eSub.Sub(&e1, &e2)
        eSubAssign.Set(&e1)
        eSubAssign.SubAssign(&e2)
        eDiv.Div(&e1, &e2)
        eNeg.Neg(&e1)
        eInv.Inverse(&e1)
		eExp.Exp(e1, rExp)
		bits := b2.Bits()
		exponent := make([]uint64, len(bits))
		for k := 0; k < len(bits); k++ {
			exponent[k] = uint64(bits[k])
		}
		eExp2.Exp(e1, exponent...)
        eLsh.Double(&e1)

        // same operations with big int
        bAdd.Add(b1, b2).Mod(&bAdd, modulus)
        bMul.Mul(b1, b2).Mod(&bMul, modulus)
        bSquare.Mul(b1, b1).Mod(&bSquare, modulus)
        bSub.Sub(b1, b2).Mod(&bSub, modulus)
        bDiv.ModInverse(b2, modulus)
        bDiv.Mul(&bDiv, b1).
            Mod(&bDiv, modulus)
        bNeg.Neg(b1).Mod(&bNeg, modulus)

        bInv.ModInverse(b1, modulus)
		bExp.Exp(b1, rbExp, modulus)
		bExp2.Exp(b1, b2, modulus)
        bLsh.Lsh(b1, 1).Mod(&bLsh, modulus)

        cmpEandB(&eSquare, &bSquare, "Square")
		cmpEandB(&eMul, &bMul, "Mul")
        cmpEandB(&eMulAssign, &bMul, "MulAssign")
        cmpEandB(&eAdd, &bAdd, "Add")
        cmpEandB(&eAddAssign, &bAdd, "AddAssign")
        cmpEandB(&eSub, &bSub, "Sub")
        cmpEandB(&eSubAssign, &bSub, "SubAssign")
        cmpEandB(&eDiv, &bDiv, "Div")
        cmpEandB(&eNeg, &bNeg, "Neg")
        cmpEandB(&eInv, &bInv, "Inv")
		cmpEandB(&eExp, &bExp, "Exp")
		cmpEandB(&eExp2, &bExp2, "Exp multi words")
		cmpEandB(&eLsh, &bLsh, "Lsh")
		
		// legendre symbol
		if e1.Legendre() != big.Jacobi(b1, modulus) {
			t.Fatal("legendre symbol computation failed")
		}
		if e2.Legendre() != big.Jacobi(b2, modulus) {
			t.Fatal("legendre symbol computation failed")
		}

		// sqrt 
		var eSqrt {{.ElementName}}
		var bSqrt big.Int
		bSqrt.ModSqrt(b1, modulus)
		eSqrt.Sqrt(&e1)
		cmpEandB(&eSqrt, &bSqrt, "Sqrt")
	}
}

func Test{{toUpper .ElementName}}IsRandom(t *testing.T) {
	for i := 0; i < 50; i++ {
		var x, y {{.ElementName}}
		x.SetRandom()
		y.SetRandom()
		if x.Equal(&y) {
			t.Fatal("2 random numbers are unlikely to be equal")
		}
	}
}

// -------------------------------------------------------------------------------------------------
// benchmarks
// most benchmarks are rudimentary and should sample a large number of random inputs
// or be run multiple times to ensure it didn't measure the fastest path of the function

var benchRes{{.ElementName}} {{.ElementName}}

func BenchmarkInverse{{toUpper .ElementName}}(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Inverse(&x)
	}

}
func BenchmarkExp{{toUpper .ElementName}}(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Exp(x, mrand.Uint64())
	}
}


func BenchmarkDouble{{toUpper .ElementName}}(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Double(&benchRes{{.ElementName}})
	}
}


func BenchmarkAdd{{toUpper .ElementName}}(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Add(&x, &benchRes{{.ElementName}})
	}
}

func BenchmarkSub{{toUpper .ElementName}}(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Sub(&x, &benchRes{{.ElementName}})
	}
}

func BenchmarkNeg{{toUpper .ElementName}}(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Neg(&benchRes{{.ElementName}})
	}
}

func BenchmarkDiv{{toUpper .ElementName}}(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Div(&x, &benchRes{{.ElementName}})
	}
}


func BenchmarkFromMont{{toUpper .ElementName}}(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.FromMont()
	}
}

func BenchmarkToMont{{toUpper .ElementName}}(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.ToMont()
	}
}
func BenchmarkSquare{{toUpper .ElementName}}(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Square(&benchRes{{.ElementName}})
	}
}

func BenchmarkSqrt{{toUpper .ElementName}}(b *testing.B) {
	var a {{.ElementName}}
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Sqrt(&a)
	}
}

func BenchmarkMulAssign{{toUpper .ElementName}}(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}}.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.MulAssign(&x)
	}
}

{{ if .ASM}}
func BenchmarkMulAsm{{toUpper .ElementName}}(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}}.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mulAsm{{.ElementName}}(&benchRes{{.ElementName}}, &x)
		// benchRes{{.ElementName}}.MulAssign(&x)
	}
}
{{ end}}


func Test{{toUpper .ElementName}}MulAsm(t *testing.T) {
	modulus, _ := new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	for i := 0; i < 1000; i++ {
		// sample 2 random big int
		b1, _ := rand.Int(rand.Reader, modulus)
		b2, _ := rand.Int(rand.Reader, modulus)

		// e1 = mont(b1), e2 = mont(b2)
		var e1, e2, eMul, eMulAsm {{.ElementName}}
		e1.SetBigInt(b1)
		e2.SetBigInt(b2)

		eMul = e1
		eMul.testMulAssign(&e2)
		eMulAsm = e1
		eMulAsm.MulAssign(&e2)

		if !eMul.Equal(&eMulAsm) {
			t.Fatal("inconsisntencies between MulAssign and testMulAssign --> check if MulAssign is calling ASM implementaiton on amd64")
		}
	}
}

// this is here for consistency purposes, to ensure MulAssign on AMD64 using asm implementation gives consistent results 
func (z *{{.ElementName}}) testMulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x"}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

{{ if .Benches}}
// Montgomery multiplication benchmarks
func (z *{{.ElementName}}) mulCIOS(x *{{.ElementName}}) *{{.ElementName}} {
	{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x"}}
	{{ template "reduce" . }}
	return z 
}

func (z *{{.ElementName}}) mulNoCarry(x *{{.ElementName}}) *{{.ElementName}} {
	{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "x"}}
	{{ template "reduce" . }}
	return z 
}

func (z *{{.ElementName}}) mulFIPS(x *{{.ElementName}}) *{{.ElementName}} {
	{{ template "mul_fips" dict "all" . "V1" "z" "V2" "x"}}
	{{ template "reduce" . }}
	return z 
}


func BenchmarkMulCIOS{{toUpper .ElementName}}(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}}.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.mulCIOS(&x)
	}
}

func BenchmarkMulFIPS{{toUpper .ElementName}}(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}}.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.mulFIPS(&x)
	}
}

func BenchmarkMulNoCarry{{toUpper .ElementName}}(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}}.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.mulNoCarry(&x)
	}
}

{{ end }}


`
