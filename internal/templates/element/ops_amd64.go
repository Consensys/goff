package element

const OpsAMD64 = `

// set functions pointers to ADX version if instruction set available
func init() {
	if supportAdx {
		{{if gt .NbWords 6}}

			mul{{.ElementName}} = func (res,x,y *{{.ElementName}}) {
				_mulLargeADX{{.ElementName}}(res{{range $i := .NbWordsIndexesFull}}, x[{{$i}}]{{end}} {{range $i := .NbWordsIndexesFull}}, y[{{$i}}]{{end}}  )
			}
			square{{.ElementName}} = func (res,x *{{.ElementName}}) {
				_mulLargeADX{{.ElementName}}(res{{range $i := .NbWordsIndexesFull}}, x[{{$i}}]{{end}} {{range $i := .NbWordsIndexesFull}}, x[{{$i}}]{{end}}  )
			}

		{{else}}
			mul{{.ElementName}} = _mulADX{{.ElementName}}
			square{{.ElementName}} = _squareADX{{.ElementName}}
		{{end}}

		fromMont{{.ElementName}} = _fromMontADX{{.ElementName}}
	}
}

// -------------------------------------------------------------------------------------------------
// Declarations

{{if gt .NbWords 6}}

//go:noescape
//go:noescape
func _mulLargeADX{{.ElementName}}(res *{{.ElementName}} {{range $i := .NbWordsIndexesFull}}, x{{$i}}{{end}} {{range $i := .NbWordsIndexesFull}}, y{{$i}}{{end}}  uint64)


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


{{else}}


//go:noescape
func _mulADX{{.ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func _squareADX{{.ElementName}}(res,x *{{.ElementName}})

//go:noescape
func sub{{.ElementName}}(res,x,y *{{.ElementName}})

// Sub  z = x - y mod q
func (z *{{.ElementName}}) Sub( x, y *{{.ElementName}}) *{{.ElementName}} {
	sub{{.ElementName}}(z, x, y)
	return z
}
// SubAssign  z = z - x mod q
func (z *{{.ElementName}}) SubAssign(x *{{.ElementName}}) *{{.ElementName}} {
	sub{{.ElementName}}(z, z, x)
	return z
}

{{end}}

//go:noescape
func reduce{{.ElementName}}(res *{{.ElementName}})

//go:noescape
func add{{.ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func double{{.ElementName}}(res,x *{{.ElementName}})

//go:noescape
func _fromMontADX{{.ElementName}}(res *{{.ElementName}})



// Add z = x + y mod q
func (z *{{.ElementName}}) Add( x, y *{{.ElementName}}) *{{.ElementName}} {
	add{{.ElementName}}(z, x, y)
	return z 
}

// AddAssign z = z + x mod q
func (z *{{.ElementName}}) AddAssign(x *{{.ElementName}}) *{{.ElementName}} {
	add{{.ElementName}}(z, z, x)
	return z 
}


// Double z = x + x mod q, aka Lsh 1
func (z *{{.ElementName}}) Double( x *{{.ElementName}}) *{{.ElementName}} {
	double{{.ElementName}}(z, x)
	return z 
}








`
