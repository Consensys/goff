package e2

// Base ...
const Base = `
import (
	"golang.org/x/sys/cpu"
)

// supportAdx will be set only on amd64 that has MULX and ADDX instructions
var supportAdx = cpu.X86.HasADX && cpu.X86.HasBMI2

// q (modulus)
var q{{.ElementName}} = [{{.NbWords}}]uint64{
	{{- range $i := .NbWordsIndexesFull}}
	{{index $.Q $i}},{{end}}
}

// q'[0], see montgommery multiplication algorithm
var q{{.ElementName}}Inv0 uint64 = {{index $.QInverse 0}}


//go:noescape
func add{{toUpper .ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func sub{{toUpper .ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func double{{toUpper .ElementName}}(res,x *{{.ElementName}})

//go:noescape
func neg{{toUpper .ElementName}}(res,x *{{.ElementName}})

{{if .BN256}}

//go:noescape
func mulNonRes{{toUpper .ElementName}}(res,x *{{.ElementName}})

//go:noescape
func squareAdx{{toUpper .ElementName}}(res,x *{{.ElementName}})

//go:noescape
func mulAdx{{toUpper .ElementName}}(res,x,y *{{.ElementName}})


// MulByNonResidue multiplies a {{.ElementName}} by (9,1)
func (z *{{.ElementName}}) MulByNonResidue(x *{{.ElementName}}) *{{.ElementName}} {
	mulNonRes{{toUpper .ElementName}}(z, x)
	return z
}

// Mul sets z to the {{.ElementName}}-product of x,y, returns z
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	mulAdx{{toUpper .ElementName}}(z, x, y)
	return z
}

// Square sets z to the {{.ElementName}}-product of x,x, returns z
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	squareAdx{{toUpper .ElementName}}(z, x)
	return z
}


{{end}}

`
