package element

// Ops is included with all builds (regardless of architecture or if F.ASM is set)
const Ops = `


// Generic (no ADX instructions, no AMD64) versions

func _mulGeneric{{.ElementName}}(z,x,y *{{.ElementName}}) {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "y"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "y" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}


func _squareGeneric{{.ElementName}}(z,x *{{.ElementName}}) {
	{{ if .NoCarry}}
		{{ template "mul_nocarry" dict "all" . "V1" "x" "V2" "x"}}
	{{ else }}
		{{ template "mul_cios" dict "all" . "V1" "x" "V2" "x" "NoReturn" true}}
	{{ end }}
	{{ template "reduce" . }}
}

func _fromMontGeneric{{.ElementName}}(z *{{.ElementName}}) {
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
