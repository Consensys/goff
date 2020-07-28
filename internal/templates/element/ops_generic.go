package element

const OpsNoAsm = `
// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

func Mul(z, x, y *{{.ElementName}}) {
	_mulGeneric(z, x, y)
}

func Square(z, x *{{.ElementName}}) {
	_squareGeneric(z,x)
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func FromMont(z *{{.ElementName}} ) {
	_fromMontGeneric(z)
}

// Add z = x + y mod q
func Add(z,  x, y *{{.ElementName}}) {
	var carry uint64
	{{$k := sub $.NbWords 1}}
	z[0], carry = bits.Add64(x[0], y[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
		{{- else}}
			z[{{$i}}], carry = bits.Add64(x[{{$i}}], y[{{$i}}], carry)
		{{- end}}
	{{- end}}
	{{- if .NoCarry}}
		z[{{$k}}], _ = bits.Add64(x[{{$k}}], y[{{$k}}], carry)
	{{- else }}
		z[{{$k}}], carry = bits.Add64(x[{{$k}}], y[{{$k}}], carry)
		// if we overflowed the last addition, z >= q
		// if z >= q, z = z - q
		if carry != 0 {
			// we overflowed, so z >= q
			z[0], carry = bits.Sub64(z[0], {{index $.Q 0}}, 0)
			{{- range $i := .NbWordsIndexesNoZero}}
				z[{{$i}}], carry = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, carry)
			{{- end}}
			return 
		}
	{{- end}}

	{{ template "reduce" .}}
}

// Double z = x + x mod q, aka Lsh 1
func Double(z,  x *{{.ElementName}}) {
	var carry uint64
	{{$k := sub $.NbWords 1}}
	z[0], carry = bits.Add64(x[0], x[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
		{{- else}}
			z[{{$i}}], carry = bits.Add64(x[{{$i}}], x[{{$i}}], carry)
		{{- end}}
	{{- end}}
	{{- if .NoCarry}}
		z[{{$k}}], _ = bits.Add64(x[{{$k}}], x[{{$k}}], carry)
	{{- else }}
		z[{{$k}}], carry = bits.Add64(x[{{$k}}], x[{{$k}}], carry)
		// if we overflowed the last addition, z >= q
		// if z >= q, z = z - q
		if carry != 0 {
			// we overflowed, so z >= q
			z[0], carry = bits.Sub64(z[0], {{index $.Q 0}}, 0)
			{{- range $i := .NbWordsIndexesNoZero}}
				z[{{$i}}], carry = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, carry)
			{{- end}}
			return 
		}
	{{- end}}

	{{ template "reduce" .}}
}


// Sub  z = x - y mod q
func Sub(z,  x, y *{{.ElementName}}) {
	var b uint64
	z[0], b = bits.Sub64(x[0], y[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		z[{{$i}}], b = bits.Sub64(x[{{$i}}], y[{{$i}}], b)
	{{- end}}
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], {{index $.Q 0}}, 0)
		{{- range $i := .NbWordsIndexesNoZero}}
			{{- if eq $i $.NbWordsLastIndex}}
				z[{{$i}}], _ = bits.Add64(z[{{$i}}], {{index $.Q $i}}, c)
			{{- else}}
				z[{{$i}}], c = bits.Add64(z[{{$i}}], {{index $.Q $i}}, c)
			{{- end}}
		{{- end}}
	}
}

// Neg z = q - x 
func Neg(z,  x *{{.ElementName}}) {
	if x.IsZero() {
		z.SetZero()
		return
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
}


{{- if eq .ASM false }}

// for test purposes
func reduce(z *{{.ElementName}})  {
	{{ template "reduce" . }}
}
{{- end}}


`
