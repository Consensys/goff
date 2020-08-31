package e2

// Base ...
const Base = `


// q (modulus)
var q{{.ElementName}} = [{{.NbWords}}]uint64{
	{{- range $i := .NbWordsIndexesFull}}
	{{index $.Q $i}},{{end}}
}

// q'[0], see montgommery multiplication algorithm
var q{{.ElementName}}Inv0 uint64 = {{index $.QInverse 0}}


//go:noescape
func add{{.ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func sub{{.ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func double{{.ElementName}}(res,x *{{.ElementName}})

//go:noescape
func neg{{.ElementName}}(res,x *{{.ElementName}})

{{if .BN256}}

//go:noescape
func mulNonRes{{.ElementName}}(res,x *{{.ElementName}})

//go:noescape
func squareAdx{{.ElementName}}(res,x *{{.ElementName}})

//go:noescape
func mulAdx{{.ElementName}}(res,x,y *{{.ElementName}})


// MulByNonResidue multiplies a E2 by (9,1)
func (z *E2) MulByNonResidue(x *E2) *E2 {
	mulNonResE2(z, x)
	return z
}


{{end}}

`
