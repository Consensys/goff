package element

const SquareCIOSNoCarryAMD64 = `

// Square{{.ElementName}} z = x * x mod q 
// calling this instead of z.Square(x) is prefered for performance critical path
//go:noescape
func Square{{.ElementName}}(res,x *{{.ElementName}})

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	// if z != x {
	// 	z.Set(x)
	// }
	// MulAssign{{.ElementName}}(z, x)
	Square{{.ElementName}}(z, x)
	return z
}

`
