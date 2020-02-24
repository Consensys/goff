package element

const MontgomeryMultiplication = `

// Mul z = x * y mod q
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "y"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "y"}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

// MulAssign z = z * x mod q
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x"}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}
`
