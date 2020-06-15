package element

const OpsAMD64 = `

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\


// -------------------------------------------------------------------------------------------------
// Constants

// q (modulus)
var q{{.ElementName}} = {{.ElementName}}{
	{{- range $i := .NbWordsIndexesFull}}
	{{index $.Q $i}},{{end}}
}

// q'[0], see montgommery multiplication algorithm
var q{{.ElementName}}Inv0 uint64 = {{index $.QInverse 0}}

// rSquare
var rSquare{{.ElementName}} = {{.ElementName}}{
	{{- range $i := .RSquare}}
	{{$i}},{{end}}
}

// -------------------------------------------------------------------------------------------------
// Declarations

//go:noescape
func mul{{.ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func add{{.ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func sub{{.ElementName}}(res,x,y *{{.ElementName}})

//go:noescape
func double{{.ElementName}}(res,y *{{.ElementName}})

//go:noescape
func fromMont{{.ElementName}}(res *{{.ElementName}}) 

//go:noescape
func reduce{{.ElementName}}(res *{{.ElementName}})  // for test purposes

//go:noescape
func square{{.ElementName}}(res,y *{{.ElementName}})


// -------------------------------------------------------------------------------------------------
// APIs

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func (z *{{.ElementName}}) FromMont() *{{.ElementName}} {
	fromMont{{.ElementName}}(z)
	return z
}
	
// ToMont converts z to Montgomery form
// sets and returns z = z * r^2
func (z *{{.ElementName}}) ToMont() *{{.ElementName}} {
	mul{{.ElementName}}(z, z, &rSquare{{.ElementName}})
	return z
}

// Mul z = x * y mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	mul{{.ElementName}}(z, x, y)
	return z
}

// MulAssign z = z * x mod q 
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	mul{{.ElementName}}(z,z,x)
	return z 
}

// Add z = x + y mod q
func (z *{{.ElementName}}) Add( x, y *{{.ElementName}}) *{{.ElementName}} {
	add{{.ElementName}}(z, x,y)
	return z
}

// AddAssign z = z + x mod q
func (z *{{.ElementName}}) AddAssign(x *{{.ElementName}}) *{{.ElementName}} {
	add{{.ElementName}}(z, z,x )
	return z 
}

// Double z = x + x mod q, aka Lsh 1
func (z *{{.ElementName}}) Double( x *{{.ElementName}}) *{{.ElementName}} {
	double{{.ElementName}}(z, x)
	return z 
}

// Sub  z = x - y mod q
func (z *{{.ElementName}}) Sub( x, y *{{.ElementName}}) *{{.ElementName}} {
	sub{{.ElementName}}(z, x,y)
	return z
}

// SubAssign  z = z - x mod q
func (z *{{.ElementName}}) SubAssign(x *{{.ElementName}}) *{{.ElementName}} {
	sub{{.ElementName}}(z,z, x)
	return z 
}

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	square{{.ElementName}}(z, x)
	return z
}

`
