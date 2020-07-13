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

// field modulus stored as big.Int 
var _{{.ElementName}}Modulus big.Int 
var once{{.ElementName}}Modulus sync.Once

// {{.ElementName}}Modulus returns q as a big.Int
// q = 
// 
// {{.Modulus}}
func {{.ElementName}}Modulus() *big.Int {
	once{{.ElementName}}Modulus.Do(func() {
		_{{.ElementName}}Modulus.SetString("{{.Modulus}}", 10)
	})
	return &_{{.ElementName}}Modulus
}

// q (modulus)
var q{{.ElementName}} = {{.ElementName}}{
	{{- range $i := .NbWordsIndexesFull}}
	{{index $.Q $i}},{{end}}
}

// q'[0], see montgommery multiplication algorithm
var q{{.ElementName}}Inv0 uint64 = {{index $.QInverse 0}}

// rSquare
var rSquare{{.ElementName}} = {{.ElementName}}{
	{{- range $i := .RSquare}}
	{{$i}},{{end}}
}


// Bytes returns the regular (non montgomery) value 
// of z as a big-endian byte slice.
func (z *{{.ElementName}}) Bytes() []byte {
	var _z {{.ElementName}}
	_z.Set(z).FromMont()
	res := make([]byte, {{.ElementName}}Limbs*8)
	binary.BigEndian.PutUint64(res[({{.ElementName}}Limbs-1)*8:], _z[0])
	for i := {{.ElementName}}Limbs - 2; i >= 0; i-- {
		binary.BigEndian.PutUint64(res[i*8:(i+1)*8], _z[{{.ElementName}}Limbs-1-i])
	}
	return res
}

// SetBytes interprets e as the bytes of a big-endian unsigned integer, 
// sets z to that value (in Montgomery form), and returns z.
func (z *{{.ElementName}}) SetBytes(e []byte) *{{.ElementName}} {
	var tmp big.Int
	tmp.SetBytes(e)
	z.SetBigInt(&tmp)
	return z
}

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

// SetInterface converts i1 from uint64, int, string, or {{.ElementName}}, big.Int into {{.ElementName}}
// panic if provided type is not supported
func (z *{{.ElementName}}) SetInterface(i1 interface{}) *{{.ElementName}} {
	switch c1 := i1.(type) {
	case {{.ElementName}}:
		return z.Set(&c1)
	case *{{.ElementName}}:
		return z.Set(c1)
	case uint64:
		return z.SetUint64(c1)
	case int:
		return z.SetString(strconv.Itoa(c1))
	case string:
		return z.SetString(c1)
	case *big.Int:
		return z.SetBigInt(c1)
	case big.Int:
		return z.SetBigInt(&c1)
	case []byte:
		return z.SetBytes(c1)
	default:
		panic("invalid type")
	}
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

// One returns 1 (in montgommery form)
func One() {{.ElementName}} {
	var one {{.ElementName}}
	one.SetOne()
	return one
}

{{end}}

// MulAssign is deprecated, use Mul instead
func (z *Element) MulAssign(x *Element) *Element {
	return z.Mul(z, x)
}

// AddAssign is deprecated, use Add instead
func (z *Element) AddAssign(x *Element) *Element {
	return z.Add(z, x)
}

// SubAssign is deprecated, use Sub instead
func (z *Element) SubAssign(x *Element) *Element {
	return z.Sub(z, x)
}



`
