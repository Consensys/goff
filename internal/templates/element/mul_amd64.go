package element

const MontgomeryMultiplicationAMD64 = `

//go:noescape
func mulAssign{{.ElementName}}(res,y *{{.ElementName}})

//go:noescape
func fromMont{{.ElementName}}(z *{{.ElementName}}) 


// Mul z = x * y mod q (constant time)
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	if z == x {
		mulAssign{{.ElementName}}(z, y)
		return z
	} else if z == y {
		mulAssign{{.ElementName}}(z, x)
		return z
	} else {
		z.Set(x)
		mulAssign{{.ElementName}}(z, y)
		return z
	}
}

// MulAssign z = z * x mod q (constant time)
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	mulAssign{{.ElementName}}(z, x)
	return z 
}
`
