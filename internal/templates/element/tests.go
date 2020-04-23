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
		n = 20
	} else {
		n = 500
	}

	sAdx := supportAdx

    for i := 0; i < n; i++ {
		if i == n/2 && sAdx {
			supportAdx = false // testing without adx instruction
		}
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
        var e1, e2, eMul,  eAdd, eSub, eDiv, eNeg, eLsh, eInv, eExp, eSquare, eMulAssign, eSubAssign, eAddAssign {{.ElementName}}
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
		
		cmpEandB(&eLsh, &bLsh, "Lsh")
		
		// legendre symbol
		if e1.Legendre() != big.Jacobi(b1, modulus) {
			t.Fatal("legendre symbol computation failed")
		}
		if e2.Legendre() != big.Jacobi(b2, modulus) {
			t.Fatal("legendre symbol computation failed")
		}

		// these are slow, killing circle ci
		if n <= 5 {
			// sqrt 
			var eSqrt, eExp2 {{.ElementName}}
			var bSqrt big.Int
			bSqrt.ModSqrt(b1, modulus)
			eSqrt.Sqrt(&e1)
			cmpEandB(&eSqrt, &bSqrt, "Sqrt")
	
			bits := b2.Bits()
			exponent := make([]uint64, len(bits))
			for k := 0; k < len(bits); k++ {
				exponent[k] = uint64(bits[k])
			}
			eExp2.Exp(e1, exponent...)
			bExp2.Exp(b1, b2, modulus)
			cmpEandB(&eExp2, &bExp2, "Exp multi words")
		}
	}
	supportAdx = sAdx
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

func TestByte{{.ElementName}}(t *testing.T) {

	modulus := {{.ElementName}}Modulus()

	// test values
	var bs [3][]byte
	r1, _ := rand.Int(rand.Reader, modulus)
	bs[0] = r1.Bytes() // should be r1 as {{.ElementName}}
	r2, _ := rand.Int(rand.Reader, modulus)
	r2.Add(modulus, r2)
	bs[1] = r2.Bytes() // should be r2 as {{.ElementName}}
	var tmp big.Int
	tmp.SetUint64(0)
	bs[2] = tmp.Bytes() // should be 0 as {{.ElementName}}

	// witness values as {{.ElementName}}
	var el [3]{{.ElementName}}
	el[0].SetBigInt(r1)
	el[1].SetBigInt(r2)
	el[2].SetUint64(0)

	// check conversions
	for i := 0; i < 3; i++ {
		var z {{.ElementName}}
		z.SetBytes(bs[i])
		if !z.Equal(&el[i]) {
			t.Fatal("SetBytes fails")
		}
		// check conversion {{.ElementName}} to Bytes
		b := z.Bytes()
		z.SetBytes(b)
		if !z.Equal(&el[i]) {
			t.Fatal("Bytes fails")
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

{{if .ASM}}
func Test{{toUpper .ElementName}}Asm(t *testing.T) {
	// ensure ASM implementations matches the ones using math/bits
	modulus, _ := new(big.Int).SetString("{{.Modulus}}", 10)
	sadx := supportAdx
	for i := 0; i < 500; i++ {
		// sample 2 random big int
		if i == 250 && sadx {
			// going the no_adx path
			supportAdx = false
		}
		b1, _ := rand.Int(rand.Reader, modulus)
		b2, _ := rand.Int(rand.Reader, modulus)

		// e1 = mont(b1), e2 = mont(b2)
		var e1, e2, eTestMul, eMulAssign, eSquare, eTestSquare {{.ElementName}}
		e1.SetBigInt(b1)
		e2.SetBigInt(b2)

		eTestMul = e1
		eTestMul.testMulAssign(&e2)
		eMulAssign = e1
		eMulAssign.MulAssign(&e2)

		if !eTestMul.Equal(&eMulAssign) {
			if supportAdx {
				t.Fatal("mul assembly implementation WITH adx instructions doesn't match non-assembly one")
			} else {
				t.Fatal("mul assembly implementation WITHOUT adx instructions doesn't match non-assembly one")
			}
		}

		// square 
		eSquare.Square(&e1)
		eTestSquare.testSquare(&e1)

		if !eTestSquare.Equal(&eSquare) {
			if supportAdx {
				t.Fatal("square assembly implementation WITH adx instructions doesn't match non-assembly one")
			} else {
				t.Fatal("square assembly implementation WITHOUT adx instructions doesn't match non-assembly one")
			}
		}
	}
	supportAdx = sadx
}

func Test{{toUpper .ElementName}}reduce(t *testing.T) {
	q := {{.ElementName}} {
		{{- range $i := .NbWordsIndexesFull}}
		{{index $.Q $i}},
		{{- end}}
	}

	var testData []{{.ElementName}}
	{
		a := q
		a[{{.NbWordsLastIndex}}] -= 1 
		testData = append(testData, a)
	}
	{
		a := q
		a[0] -= 1 
		testData = append(testData, a)
	}
	{
		a := q
		a[{{.NbWordsLastIndex}}] += 1 
		testData = append(testData, a)
	}
	{
		a := q
		a[0] += 1 
		testData = append(testData, a)
	}
	{
		a := q
		testData = append(testData, a)
	}

	for _, s := range testData {
		expected := s
		reduce{{.ElementName}}(&s)
		expected.testReduce()
		if !s.Equal(&expected) {
			t.Fatal("reduce failed")
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

func (z *{{.ElementName}}) testReduce() *{{.ElementName}} {
	{{ template "reduce" . }}
	return z 
}

// this is here for consistency purposes, to ensure Square on AMD64 using asm implementation gives consistent results 
func (z *{{.ElementName}}) testSquare(x *{{.ElementName}}) *{{.ElementName}} {
	{{if .NoCarrySquare}}
		{{ template "square" dict "all" . "V1" "x"}}
		{{ template "reduce" . }}
		return z 
	{{else if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "x"}}
		{{ template "reduce" . }}
		return z 
	{{else }}
		return z.Mul(x, x)
	{{end}}
}

{{end}}



{{ if .Benches}}
// Montgomery multiplication benchmarks
func (z *{{.ElementName}}) mulCIOS(x *{{.ElementName}}) *{{.ElementName}} {
	{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x" "NoReturn" false}}
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
