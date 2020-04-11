package element

const SquareCIOSNoCarry = `

// Square z = x * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Square(x *{{.ElementName}}) *{{.ElementName}} {
	{{if .ASM}}
		if z != x {
			z.Set(x)
		}
		mulAssign{{.ElementName}}(z, x)
		return z
	{{else}}
		{{if .NoCarrySquare}}
			{{ template "square" dict "all" . "V1" "x"}}
			{{ template "reduce" . }}
			return z 
		{{else if .NoCarry}}
			{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "x"}}
			{{ template "reduce" . }}
			return z 
		{{else }}
			return z.Mul(x, x)
		{{end}}
	{{end}}
}

`
