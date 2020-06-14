package element

const OpsNoAsm = `
// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "y"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "y" "NoReturn" false}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

// MulAssign z = z * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x" "NoReturn" false}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
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

// Add z = x + y mod q
func (z *{{.ElementName}}) Add( x, y *{{.ElementName}}) *{{.ElementName}} {
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
			return z
		}
	{{- end}}

	{{ template "reduce" .}}
	return z 
}

// AddAssign z = z + x mod q
func (z *{{.ElementName}}) AddAssign(x *{{.ElementName}}) *{{.ElementName}} {
	var carry uint64
	{{$k := sub $.NbWords 1}}
	z[0], carry = bits.Add64(z[0], x[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- if eq $i $.NbWordsLastIndex}}
		{{- else}}
			z[{{$i}}], carry = bits.Add64(z[{{$i}}], x[{{$i}}], carry)
		{{- end}}
	{{- end}}
	{{- if .NoCarry}}
		z[{{$k}}], _ = bits.Add64(z[{{$k}}], x[{{$k}}], carry)
	{{- else }}
		z[{{$k}}], carry = bits.Add64(z[{{$k}}], x[{{$k}}], carry)
		// if we overflowed the last addition, z >= q
		// if z >= q, z = z - q
		if carry != 0 {
			// we overflowed, so z >= q
			z[0], carry = bits.Sub64(z[0], {{index $.Q 0}}, 0)
			{{- range $i := .NbWordsIndexesNoZero}}
				z[{{$i}}], carry = bits.Sub64(z[{{$i}}], {{index $.Q $i}}, carry)
			{{- end}}
			return z
		}
	{{- end}}

	{{ template "reduce" .}}
	return z 
}


// Double z = x + x mod q, aka Lsh 1
func (z *{{.ElementName}}) Double( x *{{.ElementName}}) *{{.ElementName}} {
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
			return z
		}
	{{- end}}

	{{ template "reduce" .}}
	return z 
}


// Sub  z = x - y mod q
func (z *{{.ElementName}}) Sub( x, y *{{.ElementName}}) *{{.ElementName}} {
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
	return z
}

// SubAssign  z = z - x mod q
func (z *{{.ElementName}}) SubAssign(x *{{.ElementName}}) *{{.ElementName}} {
	var b uint64
	z[0], b = bits.Sub64(z[0], x[0], 0)
	{{- range $i := .NbWordsIndexesNoZero}}
		z[{{$i}}], b = bits.Sub64(z[{{$i}}], x[{{$i}}], b)
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
	return z
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *{{.ElementName}}) FromMont() *{{.ElementName}} {
	// the following lines implement z = z * 1
	// with a modified CIOS montgomery multiplication
	{{- range $j := .NbWordsIndexesFull}}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * {{index $.QInverse 0}}
		C := madd0(m, {{index $.Q 0}}, z[0])
		{{- range $i := $.NbWordsIndexesNoZero}}
			C, z[{{sub $i 1}}] = madd2(m, {{index $.Q $i}}, z[{{$i}}], C)
		{{- end}}
		z[{{sub $.NbWords 1}}] = C
	}
	{{- end}}

	{{ template "reduce" .}}
	return z
}


{{- if eq .ASM false }}
func mulAssign{{.ElementName}}(z,x *{{.ElementName}}) {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}



func square{{.ElementName}}(z,x *{{.ElementName}}) {
	{{if .NoCarrySquare}}
		{{ template "square" dict "all" . "V1" "x"}}
		{{ template "reduce" . }}
	{{else if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "x"}}
		{{ template "reduce" . }}
	{{else }}
		z.Mul(x, x)
	{{end}}
}

// for test purposes
func reduce{{.ElementName}}(z *{{.ElementName}})  {
	{{ template "reduce" . }}
}
{{- end}}


`
