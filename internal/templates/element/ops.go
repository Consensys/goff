package element

// Ops is included with all builds (regardless of architecture or if F.ASM is set)
const Ops = `

import "math/bits"

{{if .ASM}}
// -------------------------------------------------------------------------------------------------
// Declarations


//go:noescape
func Add(res,x,y *{{.ElementName}})

//go:noescape
func Sub(res,x,y *{{.ElementName}})

//go:noescape
func Neg(res,x *{{.ElementName}})

//go:noescape
func Double(res,x *{{.ElementName}})

//go:noescape
func Mul(res,x,y *{{.ElementName}})

//go:noescape
func Square(res,x *{{.ElementName}})

//go:noescape
func FromMont(res *{{.ElementName}})

//go:noescape
func Reduce(res *{{.ElementName}})


// E2


{{end}}

// Generic (no ADX instructions, no AMD64) versions

func _mulGeneric(z,x,y *{{.ElementName}}) {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "y"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "y" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}


func _squareGeneric(z,x *{{.ElementName}}) {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "x" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}

func _fromMontGeneric(z *{{.ElementName}}) {
	// the following lines implement z = z * 1
	// with a modified CIOS montgomery multiplication
	{{- range $j := .NbWordsIndexesFull}}
	{
		// m = z[0]n'[0] mod W
		m := z[0] * {{index $.QInverse 0}}
		C := madd0(m, {{index $.Q 0}}, z[0])
		{{- range $i := $.NbWordsIndexesNoZero}}
			C, z[{{sub $i 1}}] = madd2(m, {{index $.Q $i}}, z[{{$i}}], C)
		{{- end}}
		z[{{sub $.NbWords 1}}] = C
	}
	{{- end}}

	{{ template "reduce" .}}
}


`
