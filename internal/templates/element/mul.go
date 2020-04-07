package element

const MontgomeryMultiplication = `
// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular, 
// there is no security guarantees such as constant time implementation 
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

// Mul z = x * y mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) Mul(x, y *{{.ElementName}}) *{{.ElementName}} {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "y"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "y" "NoReturn" false}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

// MulAssign z = z * x mod q
// see https://hackmd.io/@zkteam/modular_multiplication
func (z *{{.ElementName}}) MulAssign(x *{{.ElementName}}) *{{.ElementName}} {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x" "NoReturn" false}}
	{{ end }}
	{{ template "reduce" . }}
	return z 
}

{{- if eq .ASM false }}
// MulAssign z = z * x mod q
func MulAssign{{.ElementName}}(z,x *{{.ElementName}}) {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "z" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "z" "V2" "x" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}
{{- end}}


`
