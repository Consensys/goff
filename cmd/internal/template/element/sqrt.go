package element

const Sqrt = `
func (z *{{.ElementName}}) Legendre() int {
	var l {{.ElementName}}
	// z^((p-1)/2)
	l.Exp(*z, {{range $i := .LegendreExponent}}
		{{$i}},{{end}}
	)
	
	if l.IsZero() {
		return 0
	} 

	// if l == 1
	if {{- range $i :=  reverse .NbWordsIndexesNoZero}}(l[{{$i}}] == {{index $.One $i}}) &&{{end}}(l[0] == {{index $.One 0}})  {
		return 1
	}
	return -1
}

// Sqrt z = √x mod q
// if the square root doesn't exist (x is not a square mod q)
// Sqrt leaves z unchanged and returns nil
func (z *{{.ElementName}}) Sqrt(x *{{.ElementName}}) *{{.ElementName}} {
	switch x.Legendre() {
	case -1:
		return nil
	case 0:
		return z.SetZero()
	case 1:
		break
	}

	{{- if .Q3Mod4}}
		// q ≡ 3 (mod 4)
		// using  z ≡ ± x^((p+1)/4) (mod q)
		return z.Exp(*x, {{range $i := .Q3Mod4SqrtExponent}}
			{{$i}},{{end}}
		)
	{{- else}}
		panic("not implemented")	
	{{- end}}
}
`
