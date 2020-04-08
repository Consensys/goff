package element

const MontgomeryMultiplicationAMD64 = `

// MulAssign{{.ElementName}} z = z * x mod q (constant time)
// calling this instead of z.MulAssign(x) is prefered for performance critical path
//go:noescape
func MulAssign{{.ElementName}}(res,y *{{.ElementName}})

// Mul z = x * y mod q (constant time)
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	res := *x
	MulAssign{{.ElementName}}(&res, y)
	z.Set(&res)
	return z
}

// MulAssign z = z * x mod q (constant time)
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	MulAssign{{.ElementName}}(z, x)
	return z 
}
`
