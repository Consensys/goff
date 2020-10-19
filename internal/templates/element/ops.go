package element

// OpsAMD64 is included with AMD64 builds (regardless of architecture or if F.ASM is set)
const OpsAMD64 = `

{{if .ASM}}

// q'[0], see montgommery multiplication algorithm
var (
	q{{.ElementName}}Inv0 uint64 = {{index $.QInverse 0}}
	_ = q{{.ElementName}}Inv0 // used in asm
)

//go:noescape
func add(res,x,y *{{.ElementName}})

//go:noescape
func sub(res,x,y *{{.ElementName}})

//go:noescape
func neg(res,x *{{.ElementName}})

//go:noescape
func double(res,x *{{.ElementName}})

//go:noescape
func mul(res,x,y *{{.ElementName}})

//go:noescape
func square(res,x *{{.ElementName}})

//go:noescape
func fromMont(res *{{.ElementName}})

//go:noescape
func reduce(res *{{.ElementName}})


{{end}}



`
