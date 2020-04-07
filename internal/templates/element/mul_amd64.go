package element

const MontgomeryMultiplicationAMD64 = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

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
