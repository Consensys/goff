package element

const Ops = `

// -------------------------------------------------------------------------------------------------
// Constants

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


// -------------------------------------------------------------------------------------------------
// declarations
// do modify tests.go with new declarations to ensure both path (ADX and generic) are tested
var mul{{.ElementName}} func (res,x,y *{{.ElementName}}) = _mulGeneric{{.ElementName}}
var square{{.ElementName}} func (res,x *{{.ElementName}}) = _squareGeneric{{.ElementName}}
var fromMont{{.ElementName}} func (res *{{.ElementName}}) = _fromMontGeneric{{.ElementName}}

// -------------------------------------------------------------------------------------------------
// APIs

// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
func (z *{{.ElementName}}) ToMont() *{{.ElementName}} {
	mul{{.ElementName}}(z, z, &rSquare{{.ElementName}})
	return z
}

// Mul z = x * y mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	mul{{.ElementName}}(z, x, y)
	return z
}

// MulAssign z = z * x mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	mul{{.ElementName}}(z,z,x)
	return z 
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	square{{.ElementName}}(z,x)
	return z
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *{{.ElementName}}) FromMont() *{{.ElementName}} {
	fromMont{{.ElementName}}(z)
	return z
}




// Generic (no ADX instructions, no AMD64) versions

func _mulGeneric{{.ElementName}}(z,x,y *{{.ElementName}}) {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "y"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "y" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}


func _squareGeneric{{.ElementName}}(z,x *{{.ElementName}}) {
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

func _fromMontGeneric{{.ElementName}}(z *{{.ElementName}}) {
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
}


`
