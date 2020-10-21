package element

const OpsNoAsm = `
// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

func mul(z, x, y *{{.ElementName}}) {
	_mulGeneric(z, x, y)
}

func square(z, x *{{.ElementName}}) {
	_squareGeneric(z,x)
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func fromMont(z *{{.ElementName}} ) {
	_fromMontGeneric(z)
}

func add(z,  x, y *{{.ElementName}}) {
	_addGeneric(z,x,y)
}

func double(z,  x *{{.ElementName}}) {
	_doubleGeneric(z,x)
}


func sub(z,  x, y *{{.ElementName}}) {
	_subGeneric(z,x,y)
}

func neg(z,  x *{{.ElementName}}) {
	_negGeneric(z,x)
}


func reduce(z *{{.ElementName}})  {
	_reduceGeneric(z)
}


`
