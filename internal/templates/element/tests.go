package element

const Test = `

import (
	"crypto/rand"
	"math/big"
	"math/bits"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)


// -------------------------------------------------------------------------------------------------
// benchmarks
// most benchmarks are rudimentary and should sample a large number of random inputs
// or be run multiple times to ensure it didn't measure the fastest path of the function

var benchRes{{.ElementName}} {{.ElementName}}

func Benchmark{{toTitle .ElementName}}Inverse(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Inverse(&x)
	}

}
func Benchmark{{toTitle .ElementName}}Exp(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b1, _ := rand.Int(rand.Reader, Modulus())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Exp(x, b1)
	}
}


func Benchmark{{toTitle .ElementName}}Double(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Double(&benchRes{{.ElementName}})
	}
}


func Benchmark{{toTitle .ElementName}}Add(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Add(&x, &benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Sub(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Sub(&x, &benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Neg(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Neg(&benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Div(b *testing.B) {
	var x {{.ElementName}}
	x.SetRandom()
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Div(&x, &benchRes{{.ElementName}})
	}
}


func Benchmark{{toTitle .ElementName}}FromMont(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.FromMont()
	}
}

func Benchmark{{toTitle .ElementName}}ToMont(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.ToMont()
	}
}
func Benchmark{{toTitle .ElementName}}Square(b *testing.B) {
	benchRes{{.ElementName}}.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Square(&benchRes{{.ElementName}})
	}
}

func Benchmark{{toTitle .ElementName}}Sqrt(b *testing.B) {
	var a {{.ElementName}}
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchRes{{.ElementName}}.Sqrt(&a)
	}
}

func Benchmark{{toTitle .ElementName}}Mul(b *testing.B) {
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




func Test{{toTitle .ElementName}}SetInterface(t *testing.T) {
	// TODO 
	t.Skip("not implemented")
}

func Test{{toTitle .ElementName}}IsRandom(t *testing.T) {
	for i := 0; i < 50; i++ {
		var x, y {{.ElementName}}
		x.SetRandom()
		y.SetRandom()
		if x.Equal(&y) {
			t.Fatal("2 random numbers are unlikely to be equal")
		}
	}
}

func Test{{toTitle .ElementName}}Bytes(t *testing.T) {

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
// Gopter tests
// most of them are generated with a template

{{ if gt .NbWords 6}}
const (
	nbFuzzShort = 20
	nbFuzz = 100
)
{{else}}
const (
	nbFuzzShort = 200
	nbFuzz = 1000
)
{{end}}

// special values to be used in tests
var staticTestValues []{{.ElementName}}

func init() {
	staticTestValues = append(staticTestValues, {{.ElementName}}{}) // zero
	staticTestValues = append(staticTestValues, One()) 				// one
	staticTestValues = append(staticTestValues, rSquare) 			// r^2
	var e, one {{.ElementName}}
	one.SetOne()
	e.Sub(&q{{.ElementName}}, &one)
	staticTestValues = append(staticTestValues, e) 	// q - 1
	e.Double(&one)
	staticTestValues = append(staticTestValues, e) 	// 2 

	{
		a := q{{.ElementName}}
		a[{{.NbWordsLastIndex}}]--
		staticTestValues = append(staticTestValues, a)
	}
	{
		a := q{{.ElementName}}
		a[0]--
		staticTestValues = append(staticTestValues, a)
	}

	{
		a := q{{.ElementName}}
		a[{{.NbWordsLastIndex}}]--
		a[0]++
		staticTestValues = append(staticTestValues, a)
	}

}


func Test{{toTitle .ElementName}}Reduce(t *testing.T) {
	testValues := make([]{{.ElementName}}, len(staticTestValues))
	copy(testValues, staticTestValues)

	for _, s := range testValues {
		expected := s
		reduce(&s)
		_reduceGeneric(&expected)
		if !s.Equal(&expected) {
			t.Fatal("reduce failed: asm and generic impl don't match")
		}
	}


	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := genFull()

	properties.Property("reduce should output a result smaller than modulus", prop.ForAll(
		func(a {{.ElementName}}) bool {
			b := a
			reduce(&a)
			_reduceGeneric(&b)
			return !a.biggerOrEqualModulus()  && a.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	// if we have ADX instruction enabled, test both path in assembly
	if supportAdx {
		t.Log("disabling ADX")
		supportAdx = false
		properties.TestingRun(t, gopter.ConsoleReporter(false))
		supportAdx = true 
	}
	
}


func Test{{toTitle .ElementName}}Legendre(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("legendre should output same result than big.Int.Jacobi", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			return a.element.Legendre() == big.Jacobi(&a.bigint, Modulus()) 
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	// if we have ADX instruction enabled, test both path in assembly
	if supportAdx {
		t.Log("disabling ADX")
		supportAdx = false
		properties.TestingRun(t, gopter.ConsoleReporter(false))
		supportAdx = true 
	}
	
}


{{template "testBinaryOp" dict "all" . "Op" "Add" "GenericOp" "_addGeneric"}}
{{template "testBinaryOp" dict "all" . "Op" "Sub" "GenericOp" "_subGeneric"}}
{{template "testBinaryOp" dict "all" . "Op" "Mul" "GenericOp" "_mulGeneric"}}
{{template "testBinaryOp" dict "all" . "Op" "Div"}}
{{template "testBinaryOp" dict "all" . "Op" "Exp"}}

{{template "testUnaryOp" dict "all" . "Op" "Square" "GenericOp" "_squareGeneric"}}
{{template "testUnaryOp" dict "all" . "Op" "Inverse"}}
{{template "testUnaryOp" dict "all" . "Op" "Sqrt"}}
{{template "testUnaryOp" dict "all" . "Op" "Double"  "GenericOp" "_doubleGeneric"}}
{{template "testUnaryOp" dict "all" . "Op" "Neg"  "GenericOp" "_negGeneric"}}

{{ define "testBinaryOp" }}

func Test{{toTitle .all.ElementName}}{{.Op}}(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}
	

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("{{.Op}}: having the receiver as operand should output the same result", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			var c, d {{.all.ElementName}}
			d.Set(&a.element)
			{{if eq .Op "Exp"}}
				c.{{.Op}}(a.element, &b.bigint)
				a.element.{{.Op}}(a.element, &b.bigint)
				b.element.{{.Op}}(d, &b.bigint)
			{{else}}
				c.{{.Op}}(&a.element, &b.element)
				a.element.{{.Op}}(&a.element, &b.element)
				b.element.{{.Op}}(&d, &b.element)
			{{end}}
			return a.element.Equal(&b.element) && a.element.Equal(&c) && b.element.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("{{.Op}}: operation result must match big.Int result", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			{
				var c {{.all.ElementName}}
				{{if eq .Op "Exp"}}
					c.{{.Op}}(a.element, &b.bigint)
				{{else}}
					c.{{.Op}}(&a.element, &b.element)
				{{end}}
				var d, e big.Int 
				{{- if eq .Op "Div"}}
					d.ModInverse(&b.bigint, Modulus())
					d.Mul(&d, &a.bigint).Mod(&d, Modulus())
				{{- else if eq .Op "Exp"}}
					d.Exp(&a.bigint, &b.bigint, Modulus())
				{{- else}}
					d.{{.Op}}(&a.bigint, &b.bigint).Mod(&d, Modulus())
				{{- end }}


				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				} 
			}

			// fixed elements
			// a is random
			// r takes special values
			testValues := make([]{{.all.ElementName}}, len(staticTestValues))
			copy(testValues, staticTestValues)

			for _, r := range testValues {
				var d, e, rb big.Int 
				r.ToBigIntRegular(&rb) 

				var c {{.all.ElementName}}
				{{- if eq .Op "Div"}}
					c.{{.Op}}(&a.element, &r)
					d.ModInverse(&rb, Modulus())
					d.Mul(&d, &a.bigint).Mod(&d, Modulus())
				{{- else if eq .Op "Exp"}}
					c.{{.Op}}(a.element, &rb)
					d.Exp(&a.bigint, &rb, Modulus())
				{{- else}}
					c.{{.Op}}(&a.element, &r)
					d.{{.Op}}(&a.bigint, &rb).Mod(&d, Modulus())
				{{- end }}

				{{if .GenericOp}}
					// checking generic impl against asm path
					var cGeneric {{.all.ElementName}}
					{{.GenericOp}}(&cGeneric, &a.element, &r)
					if !cGeneric.Equal(&c) {
						// need to give context to failing error.
						return false
					}
				{{end}}

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				} 
			}
			return true 
		},
		genA,
		genB,
	))

	properties.Property("{{.Op}}: operation result must be smaller than modulus", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			var c {{.all.ElementName}}
			{{if eq .Op "Exp"}}
				c.{{.Op}}(a.element, &b.bigint)
			{{else}}
				c.{{.Op}}(&a.element, &b.element)
			{{end}}
			return !c.biggerOrEqualModulus()
		},
		genA,
		genB,
	))

	{{if .GenericOp}}
	properties.Property("{{.Op}}: assembly implementation must be consistent with generic one", prop.ForAll(
		func(a, b testPair{{.all.ElementName}}) bool {
			var c,d {{.all.ElementName}}
			c.{{.Op}}(&a.element, &b.element)
			{{.GenericOp}}(&d, &a.element, &b.element)
			return c.Equal(&d)
		},
		genA,
		genB,
	))

	{{end}}


	specialValueTest := func() {
		// test special values against special values
		testValues := make([]{{.all.ElementName}}, len(staticTestValues))
		copy(testValues, staticTestValues)
	
		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			for _, b := range testValues {

				var bBig, d, e big.Int 
				b.ToBigIntRegular(&bBig)

				var c {{.all.ElementName}}
				


				{{- if eq .Op "Div"}}
					c.{{.Op}}(&a, &b)
					d.ModInverse(&bBig, Modulus())
					d.Mul(&d, &aBig).Mod(&d, Modulus())
				{{- else if eq .Op "Exp"}}
					c.{{.Op}}(a, &bBig)
					d.Exp(&aBig, &bBig, Modulus())
				{{- else}}
					c.{{.Op}}(&a, &b)
					d.{{.Op}}(&aBig, &bBig).Mod(&d, Modulus())
				{{- end }}
	
				{{if .GenericOp}}
					// checking asm against generic impl
					var cGeneric {{.all.ElementName}}
					{{.GenericOp}}(&cGeneric, &a, &b)
					if !cGeneric.Equal(&c) {
						t.Fatal("{{.Op}} failed special test values: asm and generic impl don't match")
					}
				{{end}}
				

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					t.Fatal("{{.Op}} failed special test values")
				} 
			}
		}
	}


	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()
	// if we have ADX instruction enabled, test both path in assembly
	if supportAdx {
		t.Log("disabling ADX")
		supportAdx = false
		properties.TestingRun(t, gopter.ConsoleReporter(false))
		specialValueTest()
		supportAdx = true 
	}
}

{{ end }}


{{ define "testUnaryOp" }}

func Test{{toTitle .all.ElementName}}{{.Op}}(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("{{.Op}}: having the receiver as operand should output the same result", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			{{if eq .Op "Sqrt"}}
			b := a.element
			{{else}}
			var b Element
			{{end}}
			b.{{.Op}}(&a.element)
			a.element.{{.Op}}(&a.element)
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("{{.Op}}: operation result must match big.Int result", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			var c {{.all.ElementName}}
			c.{{.Op}}(&a.element)

			var d, e big.Int 
			{{- if eq .Op "Square"}}
				d.Mul(&a.bigint, &a.bigint).Mod(&d, Modulus())
			{{- else if eq .Op "Inverse"}}
				d.ModInverse(&a.bigint, Modulus())
			{{- else if eq .Op "Sqrt"}}
				d.ModSqrt(&a.bigint, Modulus())
			{{- else if eq .Op "Double"}}
				d.Lsh(&a.bigint, 1).Mod(&d, Modulus())
			{{- else if eq .Op "Neg"}}
				d.Neg(&a.bigint).Mod(&d, Modulus())
			{{- end }}


			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0
		},
		genA,
	))

	properties.Property("{{.Op}}: operation result must be smaller than modulus", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			var c {{.all.ElementName}}
			c.{{.Op}}(&a.element)
			return !c.biggerOrEqualModulus()
		},
		genA,
	))

	{{if .GenericOp}}
	properties.Property("{{.Op}}: assembly implementation must be consistent with generic one", prop.ForAll(
		func(a testPair{{.all.ElementName}}) bool {
			var c,d {{.all.ElementName}}
			c.{{.Op}}(&a.element)
			{{.GenericOp}}(&d, &a.element)
			return c.Equal(&d)
		},
		genA,
	))

	{{end}}


	specialValueTest := func() {
		// test special values
		testValues := make([]{{.all.ElementName}}, len(staticTestValues))
		copy(testValues, staticTestValues)
	
		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			var c {{.all.ElementName}}
			c.{{.Op}}(&a)

			var  d, e big.Int 
			{{- if eq .Op "Square"}}
				d.Mul(&aBig, &aBig).Mod(&d, Modulus())
			{{- else if eq .Op "Inverse"}}
				d.ModInverse(&aBig, Modulus())
			{{- else if eq .Op "Sqrt"}}
				d.ModSqrt(&aBig, Modulus())
			{{- else if eq .Op "Double"}}
				d.Lsh(&aBig, 1).Mod(&d, Modulus())
			{{- else if eq .Op "Neg"}}
				d.Neg(&aBig).Mod(&d, Modulus())
			{{- end }}

			{{if .GenericOp}}
				// checking asm against generic impl
				var cGeneric {{.all.ElementName}}
				{{.GenericOp}}(&cGeneric, &a)
				if !cGeneric.Equal(&c) {
					t.Fatal("{{.Op}} failed special test values: asm and generic impl don't match")
				}
			{{end}}
			

			if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
				t.Fatal("{{.Op}} failed special test values")
			} 
		}
	}


	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()
	// if we have ADX instruction enabled, test both path in assembly
	if supportAdx {
		supportAdx = false
		t.Log("disabling ADX")
		properties.TestingRun(t, gopter.ConsoleReporter(false))
		specialValueTest()
		supportAdx = true 
	}
}

{{ end }}


func Test{{toTitle .ElementName}}FromMont(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Assembly implementation must be consistent with generic one", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			c := a.element
			d := a.element
			c.FromMont()
			_fromMontGeneric(&d)
			return c.Equal(&d)
		},
		genA,
	))

	properties.Property("x.FromMont().ToMont() == x", prop.ForAll(
		func(a testPair{{.ElementName}}) bool {
			c := a.element
			c.FromMont().ToMont()
			return c.Equal(&a.element)
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


func genFull() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {

		genRandomFq := func() {{.ElementName}} {
			var g {{.ElementName}}

			g = {{.ElementName}}{
				{{- range $i := .NbWordsIndexesFull}}
				genParams.NextUint64(),{{end}}
			}

			if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
				g[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
			}

			for g.biggerOrEqualModulus() {
				g = {{.ElementName}}{
					{{- range $i := .NbWordsIndexesFull}}
					genParams.NextUint64(),{{end}}
				}
				if q{{.ElementName}}[{{.NbWordsLastIndex}}] != ^uint64(0) {
					g[{{.NbWordsLastIndex}}] %= (q{{.ElementName}}[{{.NbWordsLastIndex}}] +1 )
				}
			}

			return g 
		}
		a := genRandomFq()

		var carry uint64
		{{- range $i := .NbWordsIndexesFull}}
			{{- if eq $i $.NbWordsLastIndex}}
			a[{{$i}}], _ = bits.Add64(a[{{$i}}], q{{$.ElementName}}[{{$i}}], carry)
			{{- else}}
			a[{{$i}}], carry = bits.Add64(a[{{$i}}], q{{$.ElementName}}[{{$i}}], carry)
			{{- end}}
		{{- end}}
		
		genResult := gopter.NewGenResult(a, gopter.NoShrinker)
		return genResult
	}
}


`
