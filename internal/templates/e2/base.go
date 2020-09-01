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

// Mul sets z to the E2-product of x,y, returns z
func (z *E2) Mul(x, y *E2) *E2 {
	mulAdxE2(z, x, y)
	return z
}

// Square sets z to the E2-product of x,x, returns z
func (z *E2) Square(x *E2) *E2 {
	squareAdxE2(z, x)
	return z
}


{{end}}

`
