package element

// OpsAMD64 generates ops_amd64.go
const OpsAMD64 = `


// -------------------------------------------------------------------------------------------------
// Declarations

//go:noescape
func reduce{{.ElementName}}(res *{{.ElementName}})

//go:noescape
func add{{.ElementName}}(res,x,y *{{.ElementName}})


//go:noescape
func double{{.ElementName}}(res,x *{{.ElementName}})

//go:noescape
func _fromMontADX{{.ElementName}}(res *{{.ElementName}})

{{if gt .NbWords 6}}

//go:noescape
func _mulLargeADX{{.ElementName}}(res *{{.ElementName}} {{range $i := .NbWordsIndexesFull}}, x{{$i}}{{end}} {{range $i := .NbWordsIndexesFull}}, y{{$i}}{{end}}  uint64)

{{else}}

//go:noescape
func _mulADX{{.ElementName}}(res,x,y *{{.ElementName}})

{{if .NoCarrySquare}}
//go:noescape
func _squareADX{{.ElementName}}(res,x *{{.ElementName}})
{{end}}
//go:noescape
func sub{{.ElementName}}(res,x,y *{{.ElementName}})

{{end}}

// -------------------------------------------------------------------------------------------------
// APIs

// Add z = x + y mod q
func (z *{{.ElementName}}) Add( x, y *{{.ElementName}}) *{{.ElementName}} {
	add{{.ElementName}}(z, x, y)
	return z 
}

// Double z = x + x mod q, aka Lsh 1
func (z *{{.ElementName}}) Double( x *{{.ElementName}}) *{{.ElementName}} {
	double{{.ElementName}}(z, x)
	return z 
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *{{.ElementName}}) FromMont() *{{.ElementName}} {
	_fromMontADX{{.ElementName}}(z)
	return z
}


{{if gt .NbWords 6}}

// Mul z = x * y mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	if supportAdx {
		_mulLargeADX{{.ElementName}}(z{{range $i := .NbWordsIndexesFull}}, x[{{$i}}]{{end}} {{range $i := .NbWordsIndexesFull}}, y[{{$i}}]{{end}}  )
	} else {
		_mulGeneric{{.ElementName}}(z, x, y)
	}
	return z
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	if supportAdx {
		_mulLargeADX{{.ElementName}}(z{{range $i := .NbWordsIndexesFull}}, x[{{$i}}]{{end}}{{range $i := .NbWordsIndexesFull}}, x[{{$i}}]{{end}} )
	} else {
		_squareGeneric{{.ElementName}}(z, x)
	}
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

{{else}}

// Sub  z = x - y mod q
func (z *{{.ElementName}}) Sub( x, y *{{.ElementName}}) *{{.ElementName}} {
	sub{{.ElementName}}(z, x, y)
	return z
}

// Mul z = x * y mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	_mulADX{{.ElementName}}(z, x, y)
	return z
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	{{if .NoCarrySquare}}
		_squareADX{{.ElementName}}(z,x)
	{{else}}
		_mulADX{{.ElementName}}(z, x, x)
	{{end}}
	return z
}


{{end}}
`
