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
	{{if eq .NoCollidingNames false}}"strconv"{{- end}}
)

// {{.ElementName}} represents a field element stored on {{.NbWords}} words (uint64)
// {{.ElementName}} are assumed to be in Montgomery form in all methods
// field modulus q =
// 
// {{.Modulus}}
type {{.ElementName}} [{{.NbWords}}]uint64

// Limbs number of 64 bits words needed to represent {{.ElementName}}
const Limbs = {{.NbWords}}

// Bits number bits needed to represent {{.ElementName}}
const Bits = {{.NbBits}}

// field modulus stored as big.Int 
var _modulus big.Int 
var onceModulus sync.Once

// Modulus returns q as a big.Int
// q = 
// 
// {{.Modulus}}
func Modulus() *big.Int {
	onceModulus.Do(func() {
		_modulus.SetString("{{.Modulus}}", 10)
	})
	return new(big.Int).Set(&_modulus)
}

// q (modulus)
var q{{.ElementName}} = {{.ElementName}}{
	{{- range $i := .NbWordsIndexesFull}}
	{{index $.Q $i}},{{end}}
}

// q'[0], see montgommery multiplication algorithm
var q{{.ElementName}}Inv0 uint64 = {{index $.QInverse 0}}

// rSquare
var rSquare = {{.ElementName}}{
	{{- range $i := .RSquare}}
	{{$i}},{{end}}
}


// Bytes returns the regular (non montgomery) value 
// of z as a big-endian byte slice.
func (z *{{.ElementName}}) Bytes() []byte {
	_z := z.ToRegular()
	var res [Limbs*8]byte
	{{- range $i := reverse .NbWordsIndexesFull}}
		{{- $j := mul $i 8}}
		{{- $k := sub $.NbWords 1}}
		{{- $k := sub $k $i}}
		{{- $jj := add $j 8}}
		binary.BigEndian.PutUint64(res[{{$j}}:{{$jj}}], _z[{{$k}}])
	{{- end}}

	return res[:]
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
	*z = {{.ElementName}}{v}
	return z.Mul(z, &rSquare) // z.ToMont()
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

// One returns 1 (in montgommery form)
func One() {{.ElementName}} {
	var one {{.ElementName}}
	one.SetOne()
	return one
}


// MulAssign is deprecated
// Deprecated: use Mul instead
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	return z.Mul(z, x)
}

// AddAssign is deprecated
// Deprecated: use Add instead
func (z *{{.ElementName}}) AddAssign(x *{{.ElementName}}) *{{.ElementName}} {
	return z.Add(z, x)
}

// SubAssign is deprecated
// Deprecated: use Sub instead
func (z *{{.ElementName}}) SubAssign(x *{{.ElementName}}) *{{.ElementName}} {
	return z.Sub(z, x)
}


// API with assembly impl

// Mul z = x * y mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	mul(z, x, y)
	return z
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	square(z,x)
	return z
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *{{.ElementName}}) FromMont() *{{.ElementName}} {
	fromMont(z)
	return z
}

// Add z = x + y mod q
func (z *{{.ElementName}}) Add( x, y *{{.ElementName}}) *{{.ElementName}} {
	add(z, x, y)
	return z 
}

// Double z = x + x mod q, aka Lsh 1
func (z *{{.ElementName}}) Double( x *{{.ElementName}}) *{{.ElementName}} {
	double(z, x)
	return z 
}


// Sub  z = x - y mod q
func (z *{{.ElementName}}) Sub( x, y *{{.ElementName}}) *{{.ElementName}} {
	sub(z, x, y)
	return z
}

// Neg z = q - x 
func (z *{{.ElementName}}) Neg( x *{{.ElementName}}) *{{.ElementName}} {
	neg(z, x)
	return z
}




`
