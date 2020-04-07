package element

const Base = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import (
	"math/big"
	"math/bits"
	"crypto/rand"
	"encoding/binary"
	"io"
	"sync"
	"unsafe"
	{{if eq .NoCollidingNames false}}"strconv"{{- end}}
)

// {{.ElementName}} represents a field element stored on {{.NbWords}} words (uint64)
// {{.ElementName}} are assumed to be in Montgomery form in all methods
// field modulus q =
// 
// {{.Modulus}}
type {{.ElementName}} [{{.NbWords}}]uint64

// {{.ElementName}}Limbs number of 64 bits words needed to represent {{.ElementName}}
const {{.ElementName}}Limbs = {{.NbWords}}

// {{.ElementName}}Bits number bits needed to represent {{.ElementName}}
const {{.ElementName}}Bits = {{.NbBits}}

var support_adx_{{.ElementName}} = cpu.X86.HasADX && cpu.X86.HasBMI2

// SetUint64 z = v, sets z LSB to v (non-Montgomery form) and convert z to Montgomery form
func (z *{{.ElementName}}) SetUint64(v uint64) *{{.ElementName}} {
	z[0] = v
	{{- range $i := .NbWordsIndexesNoZero}}
		z[{{$i}}] = 0
	{{- end}}
	return z.ToMont()
}

// Set z = x
func (z *{{.ElementName}}) Set(x *{{.ElementName}}) *{{.ElementName}} {
	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = x[{{$i}}]
	{{- end}}
	return z
}

// SetZero z = 0
func (z *{{.ElementName}}) SetZero() *{{.ElementName}} {
	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = 0
	{{- end}}
	return z
}

// SetOne z = 1 (in Montgomery form)
func (z *{{.ElementName}}) SetOne() *{{.ElementName}} {
	{{- range $i := .NbWordsIndexesFull}}
		z[{{$i}}] = {{index $.One $i}}
	{{- end}}
	return z
}


// Neg z = q - x 
func (z *{{.ElementName}}) Neg( x *{{.ElementName}}) *{{.ElementName}} {
	if x.IsZero() {
		return z.SetZero()
	}
	var borrow uint64
	z[0], borrow = bits.Sub64({{index $.Q 0}}, x[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
			z[{{$i}}], _ = bits.Sub64({{index $.Q $i}}, x[{{$i}}], borrow)
		{{- else}}
			z[{{$i}}], borrow = bits.Sub64({{index $.Q $i}}, x[{{$i}}], borrow)
		{{- end}}
	{{- end}}
	return z
}


// Div z = x*y^-1 mod q 
func (z *{{.ElementName}}) Div( x, y *{{.ElementName}}) *{{.ElementName}} {
	var yInv {{.ElementName}}
	yInv.Inverse( y)
	z.Mul( x, &yInv)
	return z
}

// Equal returns z == x
func (z *{{.ElementName}}) Equal(x *{{.ElementName}}) bool {
	return {{- range $i :=  reverse .NbWordsIndexesNoZero}}(z[{{$i}}] == x[{{$i}}]) &&{{end}}(z[0] == x[0])
}

// IsZero returns z == 0
func (z *{{.ElementName}}) IsZero() bool {
	return ( {{- range $i :=  reverse .NbWordsIndexesNoZero}} z[{{$i}}] | {{end}}z[0]) == 0
}



// field modulus stored as big.Int 
var _{{toLower .ElementName}}ModulusBigInt big.Int 
var once{{toLower .ElementName}}Modulus sync.Once
func {{toLower .ElementName}}ModulusBigInt() *big.Int {
	once{{toLower .ElementName}}Modulus.Do(func() {
		_{{toLower .ElementName}}ModulusBigInt.SetString("{{.Modulus}}", 10)
	})
	return &_{{toLower .ElementName}}ModulusBigInt
}


{{/* We use big.Int for Inverse for these type of moduli */}}
{{if eq .NoCarry false}}

// Inverse z = x^-1 mod q 
// note: allocates a big.Int (math/big)
func (z *{{.ElementName}}) Inverse( x *{{.ElementName}}) *{{.ElementName}} {
	var _xNonMont big.Int
	x.ToBigIntRegular( &_xNonMont)
	_xNonMont.ModInverse(&_xNonMont, {{toLower .ElementName}}ModulusBigInt())
	z.SetBigInt(&_xNonMont)
	return z
}

{{ else }}

// Inverse z = x^-1 mod q 
// Algorithm 16 in "Efficient Software-Implementation of Finite Fields with Applications to Cryptography"
// if x == 0, sets and returns z = x 
func (z *{{.ElementName}}) Inverse(x *{{.ElementName}}) *{{.ElementName}} {
	if x.IsZero() {
		return z.Set(x)
	}

	// initialize u = q
	var u = {{.ElementName}}{
		{{- range $i := .NbWordsIndexesFull}}
		{{index $.Q $i}},{{end}}
	}

	// initialize s = r^2
	var s = {{.ElementName}}{
		{{- range $i := .RSquare}}
		{{$i}},{{end}}
	}

	// r = 0
	r := {{.ElementName}}{}

	v := *x

	var carry, borrow, t, t2 uint64
	var bigger, uIsOne, vIsOne bool

	for !uIsOne && !vIsOne {
		for v[0]&1 == 0 {
			{{ template "div2" dict "all" . "V" "v"}}
			if s[0]&1 == 1 {
				{{ template "add_q" dict "all" . "V1" "s" }}
			}
			{{ template "div2" dict "all" . "V" "s"}}
		} 
		for u[0]&1 == 0 {
			{{ template "div2" dict "all" . "V" "u"}}
			if r[0]&1 == 1 {
				{{ template "add_q" dict "all" . "V1" "r" }}
			}
			{{ template "div2" dict "all" . "V" "r"}}
		} 
		{{ template "bigger" dict "all" . "V1" "v" "V2" "u"}}
		if bigger  {
			{{ template "sub_noborrow" dict "all" . "V1" "v" "V2" "u" }}
			{{ template "bigger" dict "all" . "V1" "r" "V2" "s"}}
			if bigger {
				{{ template "add_q" dict "all" . "V1" "s" }}
			}
			{{ template "sub_noborrow" dict "all" . "V1" "s" "V2" "r" }}
		} else {
			{{ template "sub_noborrow" dict "all" . "V1" "u" "V2" "v" }}
			{{ template "bigger" dict "all" . "V1" "s" "V2" "r"}}
			if bigger {
				{{ template "add_q" dict "all" . "V1" "r" }}
			}
			{{ template "sub_noborrow" dict "all" . "V1" "r" "V2" "s" }}
		}
		uIsOne = (u[0] == 1) && ({{- range $i := reverse .NbWordsIndexesNoZero}}u[{{$i}}] {{if eq $i 1}}{{else}} | {{end}}{{end}} ) == 0
		vIsOne = (v[0] == 1) && ({{- range $i := reverse .NbWordsIndexesNoZero}}v[{{$i}}] {{if eq $i 1}}{{else}} | {{end}}{{end}} ) == 0
	}

	if uIsOne {
		z.Set(&r)
	} else {
		z.Set(&s)
	}

	return z
}

{{ end }}

// SetRandom sets z to a random element < q
func (z *{{.ElementName}}) SetRandom() *{{.ElementName}} {
	bytes := make([]byte, {{mul 8 .NbWords}})
	io.ReadFull(rand.Reader, bytes)
	{{- range $i :=  .NbWordsIndexesFull}}
		{{- $k := add $i 1}}
		z[{{$i}}] = binary.BigEndian.Uint64(bytes[{{mul $i 8}}:{{mul $k 8}}]) 
	{{- end}}
	z[{{$.NbWordsLastIndex}}] %= {{index $.Q $.NbWordsLastIndex}}

	{{ template "reduce" . }}

	return z
}

{{ if .NoCollidingNames}}
{{ else}}
func One() {{.ElementName}} {
	var one {{.ElementName}}
	one.SetOne()
	return one
}

func FromInterface(i1 interface{}) {{.ElementName}} {
	var val {{.ElementName}}

	switch c1 := i1.(type) {
	case uint64:
		val.SetUint64(c1)
	case int:
		val.SetString(strconv.Itoa(c1))
	case string:
		val.SetString(c1)
	case {{.ElementName}}:
		val = c1
	case *{{.ElementName}}:
		val.Set(c1)
	// TODO add big.Int convertions
	default:
		panic("invalid type")
	}

	return val
}
{{end}}


{{ define "bigger" }}
	// {{$.V1}} >= {{$.V2}}
	bigger = !({{- range $i := reverse $.all.NbWordsIndexesNoZero}} {{$.V1}}[{{$i}}] < {{$.V2}}[{{$i}}] || ( {{$.V1}}[{{$i}}] == {{$.V2}}[{{$i}}] && (
		{{- end}}{{$.V1}}[0] < {{$.V2}}[0] {{- range $i :=  $.all.NbWordsIndexesNoZero}} )) {{- end}} )
{{ end }}

{{ define "add_q" }}
	// {{$.V1}} = {{$.V1}} + q 
	{{$.V1}}[0], carry = bits.Add64({{$.V1}}[0], {{index $.all.Q 0}}, 0)
	{{- range $i := .all.NbWordsIndexesNoZero}}
		{{- if eq $i $.all.NbWordsLastIndex}}
			{{$.V1}}[{{$i}}], _ = bits.Add64({{$.V1}}[{{$i}}], {{index $.all.Q $i}}, carry)
		{{- else}}
			{{$.V1}}[{{$i}}], carry = bits.Add64({{$.V1}}[{{$i}}], {{index $.all.Q $i}}, carry)
		{{- end}}
	{{- end}}
{{ end }}

{{ define "sub_noborrow" }}
	// {{$.V1}} = {{$.V1}} - {{$.V2}}
	{{$.V1}}[0], borrow = bits.Sub64({{$.V1}}[0], {{$.V2}}[0], 0)
	{{- range $i := .all.NbWordsIndexesNoZero}}
		{{- if eq $i $.all.NbWordsLastIndex}}
			{{$.V1}}[{{$i}}], _ = bits.Sub64({{$.V1}}[{{$i}}], {{$.V2}}[{{$i}}], borrow)
		{{- else}}
			{{$.V1}}[{{$i}}], borrow = bits.Sub64({{$.V1}}[{{$i}}], {{$.V2}}[{{$i}}], borrow)
		{{- end}}
	{{- end}}
{{ end }}


{{ define "div2" }}
	// {{$.V}} = {{$.V}} >> 1
	{{- range $i :=  reverse .all.NbWordsIndexesNoZero}}
		{{- if eq $i $.all.NbWordsLastIndex}}
			t2 = {{$.V}}[{{$i}}] << 63
			{{$.V}}[{{$i}}] >>= 1
		{{- else}}
			t2 = {{$.V}}[{{$i}}] << 63
			{{$.V}}[{{$i}}] = ({{$.V}}[{{$i}}] >> 1) | t
		{{- end}}
		t = t2
	{{- end}}
	{{$.V}}[0] = ({{$.V}}[0] >> 1) | t
{{ end }}


`
