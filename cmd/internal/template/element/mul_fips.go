package element

const MulFIPS = `
{{ define "mul_fips" }}
	var p [{{.all.NbWords}}]uint64
	var t, u, v uint64

	{{- range $i := .all.NbWordsIndexesFull}}
		{{- range $j := $.all.NbWordsIndexesFull}}
			{{- if lt $j $i}}
				{{- $k := sub $i $j}}
				{{- if eq $j 0}}
					t, u, v = madd({{$.V1}}[{{$j}}], {{$.V2}}[{{$k}}], 0, u, v)
				{{- else}}
					t, u, v = madd({{$.V1}}[{{$j}}], {{$.V2}}[{{$k}}], t, u, v)
				{{- end}}
				t, u, v = madd(p[{{$j}}], {{index $.all.Q $k}}, t, u, v)
			{{- end}}
		{{- end}}
		{{- if eq $i 0}}
			u, v = bits.Mul64({{$.V1}}[{{$i}}], {{$.V2}}[0])
			p[{{$i}}] = v * {{index $.all.QInverse 0}}
			u, v, _ = madd(p[{{$i}}], {{index $.all.Q 0}}, 0, u, v)
		{{- else}}
			t, u, v = madd({{$.V1}}[{{$i}}], {{$.V2}}[0], t, u, v)
			p[{{$i}}] = v * {{index $.all.QInverse 0}}
			u, v, _ = madd(p[{{$i}}], {{index $.all.Q 0}}, t, u, v)
		{{- end}}
	{{- end}}
	{{- range $i := .all.IdxFIPS}}
		{{- $l := sub $i $.all.NbWords}}
		{{- $m := add $l 1}}
		{{- range $j := $.all.NbWordsIndexesFull}}
			{{- if ge $j $m}}
			{{- $k := sub $i $j}}
			{{- if eq $j $.all.NbWordsLastIndex}}
				t, u, v = madd({{$.V1}}[{{$j}}], {{$.V2}}[{{$k}}], t, u, v)
				u, v, p[{{$l}}] = madd(p[{{$j}}], {{index $.all.Q $k}}, t, u, v)
			{{- else}}
				{{- if eq $j $m}}
					t, u, v = madd({{$.V1}}[{{$j}}], {{$.V2}}[{{$k}}], 0, u, v)
				{{- else}}
					t, u, v = madd({{$.V1}}[{{$j}}], {{$.V2}}[{{$k}}], t, u, v)
				{{- end}}
				t, u, v = madd(p[{{$j}}], {{index $.all.Q $k}}, t, u, v)
			{{- end}}
			{{- end}}
		{{- end}}
	{{- end}}

	p[{{sub .all.NbWords 1}}] = v
	{{- range $i := reverse .all.NbWordsIndexesFull}}
		z[{{$i}}] = p[{{$i}}]
	{{- end}}
	// copy(z[:], p[:])

{{ end }}
`
