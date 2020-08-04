package element

const Test = `

import (
    "crypto/rand"
	"math/big"
	"math/bits"
	"testing"
	"fmt"
	mrand "math/rand"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func Test{{toUpper .ElementName}}CorrectnessAgainstBigInt(t *testing.T) {
    modulus := Modulus()
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
        // sample 3 random big int
        b1, _ := rand.Int(rand.Reader, modulus)
		b2, _ := rand.Int(rand.Reader, modulus)
		b3, _ := rand.Int(rand.Reader, modulus) // exponent
        

        // adding edge cases
        // TODO need more edge cases
        switch i {
        case 0:
			b3.SetUint64(0)
            b1.SetUint64(0)
        case 1:
            b2.SetUint64(0)
        case 2:
            b1.SetUint64(0)
            b2.SetUint64(0)
        case 3:
            b3.SetUint64(0)
        case 4:
            b3.SetUint64(1)
		case 5:
			b3.SetUint64(^uint64(0))
		case 6:
			b3.SetUint64(2)
			b1.Set(&modulusMinusOne)
		case 7:
			b2.Set(&modulusMinusOne)
		case 8:
			b1.Set(&modulusMinusOne)
			b2.Set(&modulusMinusOne)
        }


        var bMul, bAdd, bSub, bDiv, bNeg, bLsh, bInv, bExp, bSquare big.Int

        // e1 = mont(b1), e2 = mont(b2)
        var e1, e2, eMul,  eAdd, eSub, eDiv, eNeg, eLsh, eInv, eExp, eSquare {{.ElementName}}
        e1.SetBigInt(b1)
        e2.SetBigInt(b2)

        // (e1*e2).FromMont() === b1*b2 mod q ... etc
        eSquare.Square(&e1)
		eMul.Mul(&e1, &e2)
        eAdd.Add(&e1, &e2)
        eSub.Sub(&e1, &e2)
        eDiv.Div(&e1, &e2)
        eNeg.Neg(&e1)
        eInv.Inverse(&e1)
		eExp.Exp(e1, b3)
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
		bExp.Exp(b1, b3, modulus)
        bLsh.Lsh(b1, 1).Mod(&bLsh, modulus)

        cmpEandB(&eSquare, &bSquare, "Square")
		cmpEandB(&eMul, &bMul, "Mul")
        cmpEandB(&eAdd, &bAdd, "Add")
        cmpEandB(&eSub, &bSub, "Sub")
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
		if n <= 10 {
			// sqrt 
			var eSqrt {{.ElementName}}
			var bSqrt big.Int
			bSqrt.ModSqrt(b1, modulus)
			eSqrt.Sqrt(&e1)
			cmpEandB(&eSqrt, &bSqrt, "Sqrt")
		}
	}
	supportAdx = sAdx
}

func Test{{toUpper .ElementName}}SetInterface(t *testing.T) {
	// TODO 
	t.Skip("not implemented")
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

func TestByte(t *testing.T) {

	modulus := Modulus()

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
	b1, _ := rand.Int(rand.Reader, Modulus())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Exp(x, b1)
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

func BenchmarkMul{{toUpper .ElementName}}(b *testing.B) {
	x := {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}
	benchRes{{.ElementName}}.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Mul(&benchRes{{.ElementName}}, &x)
	}
}

{{if .ASM}}

func Test{{toUpper .ElementName}}reduce(t *testing.T) {
	q := {{.ElementName}} {
		{{- range $i := .NbWordsIndexesFull}}
		{{index $.Q $i}},
		{{- end}}
	}

	var testData []{{.ElementName}}
	{
		a := q
		a[{{.NbWordsLastIndex}}]--
		testData = append(testData, a)
	}
	{
		a := q
		a[0]--
		testData = append(testData, a)
	}
	{
		a := q
		a[{{.NbWordsLastIndex}}]++
		testData = append(testData, a)
	}
	{
		a := q
		a[0]++
		testData = append(testData, a)
	}
	{
		a := q
		testData = append(testData, a)
	}

	for _, s := range testData {
		expected := s
		reduce(&s)
		expected.testReduce()
		if !s.Equal(&expected) {
			t.Fatal("reduce failed")
		}
	}
	
}

func (z *{{.ElementName}}) testReduce() *{{.ElementName}} {
	{{ template "reduce" . }}
	return z 
}


{{end}}


// -------------------------------------------------------------------------------------------------
// Gopter tests

func Test{{toUpper .ElementName}}Mul(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10000

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("Having the receiver as operand should output the same result", prop.ForAll(
		func(a, b testPair{{.ElementName}}) bool {
			var c, d {{.ElementName}}
			d.Set(&a.element)
			c.Mul(&a.element, &b.element)
			a.element.Mul(&a.element, &b.element)
			b.element.Mul(&d, &b.element)
			return a.element.Equal(&b.element) && a.element.Equal(&c) && b.element.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Operation result must match big.Int result", prop.ForAll(
		func(a, b testPair{{.ElementName}}) bool {
			var c {{.ElementName}}
			c.Mul(&a.element, &b.element)

			var d, e big.Int 
			d.Mul(&a.bigint, &b.bigint).Mod(&d, Modulus())

			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0 
		},
		genA,
		genB,
	))

	properties.Property("Operation result must be smaller than modulus", prop.ForAll(
		func(a, b testPair{{.ElementName}}) bool {
			var c {{.ElementName}}
			c.Mul(&a.element, &b.element)
			return !c.biggerOrEqualModulus()
		},
		genA,
		genB,
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}



func Test{{toUpper .ElementName}}Square(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10000

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Having the receiver as operand should output the same result", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var b {{.ElementName}}
			b.Square(&a.element)
			a.element.Square(&a.element)
			return a.element.Equal(&b) 
		},
		genA,
	))

	properties.Property("Operation result must match big.Int result", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var b {{.ElementName}}
			b.Square(&a.element)

			var d, e big.Int 
			d.Mul(&a.bigint, &a.bigint).Mod(&d, Modulus())

			return b.FromMont().ToBigInt(&e).Cmp(&d) == 0 
		},
		genA,
	))

	properties.Property("Operation result must be smaller than modulus", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var b {{.ElementName}}
			b.Square(&a.element)
			return !b.biggerOrEqualModulus()
		},
		genA,
	))

	properties.Property("Square(x) == Mul(x,x)", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			var b,c  {{.ElementName}}
			b.Square(&a.element)
			c.Mul(&a.element, &a.element)
			return c.Equal(&b)
		},
		genA,
	))


	properties.TestingRun(t, gopter.ConsoleReporter(false))
}




type testPair{{.ElementName}} struct {
	element {{.ElementName}}
	bigint       big.Int
}

func (z *{{.ElementName}}) biggerOrEqualModulus() bool {
	{{- range $i :=  reverse .NbWordsIndexesNoZero}}
	if z[{{$i}}] > q{{$.ElementName}}[{{$i}}] {
		return true
	}
	if z[{{$i}}] < q{{$.ElementName}}[{{$i}}] {
		return false
	}
	{{end}}
	
	return z[0] >= q{{.ElementName}}[0]
}

func gen() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var g testPair{{.ElementName}}

		g.element = {{.ElementName}}{
			{{- range $i := .NbWordsIndexesFull}}
			genParams.NextUint64(),{{end}}
		}
		if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
			g.element[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
		}
		

		for g.element.biggerOrEqualModulus() {
			g.element = {{.ElementName}}{
				{{- range $i := .NbWordsIndexesFull}}
				genParams.NextUint64(),{{end}}
			}
			if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
				g.element[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
			}
		}

		g.element.ToBigIntRegular(&g.bigint)
		genResult := gopter.NewGenResult(g, gopter.NoShrinker)
		return genResult
	}
}


`
