package element

const SquareNoCarryTemplate = `
{{ define "square" }}
	var p [{{.all.NbWords}}]uint64
	
	var u, v uint64
	
	{{- range $j := .all.NbWordsIndexesFull}}
	{
		// round {{$j}}
		{{- $first_j_loop := eq $j 0}}
		{{- if $first_j_loop}}
			u, p[0] = bits.Mul64({{$.V1}}[0], {{$.V1}}[0])
		{{- end}}
		m := p[0] * {{index $.all.QInverse 0}}
		C := madd0(m, {{index $.all.Q 0}}, p[0])
		{{- if eq $j $.all.NbWordsLastIndex}}
			{{- range $i := $.all.NbWordsIndexesNoZero}}
				{{- if eq $i $.all.NbWordsLastIndex}}
					u, v = madd1({{$.V1}}[{{$j}}], {{$.V1}}[{{$j}}], p[{{$j}}])
					z[{{sub $.all.NbWords 1}}], z[{{sub $i 1}}]  = madd3(m, {{index $.all.Q $i}}, v, C, u)
				{{- else}}
					C, z[{{sub $i 1}}] = madd2(m, {{index $.all.Q $i}}, p[{{$i}}], C)
				{{- end}}
			{{- end}}
		{{- else}}
			{{- $firstDoubling := true}}
			{{- range $i := $.all.NbWordsIndexesNoZero}}
				{{- $last_i_loop := eq $i $.all.NbWordsLastIndex}}
				{{- $doubling_round := gt $i $j}}
				{{- if and $last_i_loop $doubling_round}} 
						{{- if and $first_j_loop $firstDoubling}}	
								_, u, v = madd1sb({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], u) 
						{{- else if $first_j_loop}}
								_ , u, v = madd1s({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], t, u)
						{{- else if $firstDoubling}}	
							_ , u, v = madd2sb({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], p[{{$i}}], u)
						{{- else}}
							_ , u, v = madd2s({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], p[{{$i}}], t, u)
						{{- end}}	
						p[{{sub $.all.NbWords 1}}], p[{{sub $i 1}}]  = madd3(m, {{index $.all.Q $i}}, v, C, u)
						{{- $firstDoubling = false}}
				{{- else if $last_i_loop}}
					p[{{sub $.all.NbWords 1}}], p[{{sub $i 1}}]  = madd3(m, {{index $.all.Q $i}}, p[{{$i}}], C, u)
				{{- else if $doubling_round}}
					{{- if and $first_j_loop $firstDoubling}}	
							var t uint64
							t, u, v = madd1sb({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], u)  
					{{- else if $first_j_loop}}
							t , u, v = madd1s({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], t, u)
					{{- else if $firstDoubling}}	
							var t uint64
							t , u, v = madd2sb({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], p[{{$i}}], u)
					{{- else}}
							t , u, v = madd2s({{$.V1}}[{{$j}}], {{$.V1}}[{{$i}}], p[{{$i}}], t, u)
					{{- end}}	
					C, p[{{sub $i 1}}] = madd2(m, {{index $.all.Q $i}}, v, C)
					{{- $firstDoubling = false}}
				{{- else if eq $j $i}}
					u, v = madd1({{$.V1}}[{{$j}}], {{$.V1}}[{{$j}}], p[{{$j}}])
					C, p[{{sub $i 1}}] = madd2(m, {{index $.all.Q $i}}, v, C)
				{{- else }}
					C, p[{{sub $i 1}}] = madd2(m, {{index $.all.Q $i}}, p[{{$i}}], C)
			{{- end}}
		{{- end}}
		{{- end}}
	}
	{{- end}}

{{ end }}
`
