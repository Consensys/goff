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

{{end}}

`
