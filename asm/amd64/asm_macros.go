// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package amd64

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/consensys/bavard/amd64"
)

// LabelRegisters write comment with friendler name to registers
func (f *FFAmd64) LabelRegisters(name string, r ...amd64.Register) {
	for i := 0; i < len(r); i++ {
		f.Comment(fmt.Sprintf("%s[%d] = %s", name, i, string(r[i])))
	}
	f.WriteLn("")
}

func (f *FFAmd64) copyElement(a, b []amd64.Register) {
	const tmpl = `// {{- range $i, $b := .B}}{{$b}}{{- if ne $.Last $i}},{{- else}} = {{ end}}{{- end}}
	{{- range $i, $a := .A}}{{$a}}{{- if ne $.Last $i}},{{- end}}{{- end}}
	COPY({{- range $i, $a := .A}}{{$a}},{{- end}}
		{{- range $i, $b := .B}}{{$b}}{{- if ne $.Last $i}},{{- end}}{{- end}})`
	var buf bytes.Buffer
	err := template.Must(template.New("").
		Parse(tmpl)).Execute(&buf, struct {
		A, B []amd64.Register
		Last int
	}{a, b, len(b) - 1})
	if err != nil {
		panic(err)
	}

	f.WriteLn(buf.String())
}

func (f *FFAmd64) reduceElement(r amd64.Register, a, b []amd64.Register) {
	const tmpl = `// reduce element({{- range $i, $a := .A}}{{$a}}{{- if ne $.Last $i}},{{ end}}{{- end}}) stores at {{.R}}
	REDUCE_AND_STORE({{.R}},{{- range $i, $a := .A}}{{$a}},{{- end}}
		{{- range $i, $b := .B}}{{$b}}{{- if ne $.Last $i}},{{- end}}{{- end}})`
	var buf bytes.Buffer
	err := template.Must(template.New("").
		Parse(tmpl)).Execute(&buf, struct {
		R    amd64.Register
		A, B []amd64.Register
		Last int
	}{r, a, b, len(b) - 1})
	if err != nil {
		panic(err)
	}

	f.WriteLn(buf.String())
	f.WriteLn("")
}

const tmplDefines = `

// modulus q
{{- range $i, $w := .Q}}
DATA q<>+{{mul $i 8}}(SB)/8, {{imm $w}}
{{- end}}
GLOBL q<>(SB), (RODATA+NOPTR), ${{mul 8 $.NbWords}}

// qInv0 q'[0]
DATA qInv0<>(SB)/8, {{$qinv0 := index .QInverse 0}}{{imm $qinv0}}
GLOBL qInv0<>(SB), (RODATA+NOPTR), $8

// COPY b = a 
#define COPY(ra0{{range $i := .NbWordsIndexesNoZero}},ra{{$i}}{{end}},rb0{{range $i := .NbWordsIndexesNoZero}},rb{{$i}}{{end}}) \
	{{- range $i := .NbWordsIndexesFull}}
	MOVQ ra{{$i}}, rb{{$i}};  \
	{{- end}}

// REDUCE_AND_STORE
#define REDUCE_AND_STORE(m0, ra0{{range $i := .NbWordsIndexesNoZero}},ra{{$i}}{{end}},rb0{{range $i := .NbWordsIndexesNoZero}},rb{{$i}}{{end}}) \
	COPY(ra0{{range $i := .NbWordsIndexesNoZero}},ra{{$i}}{{end}},rb0{{range $i := .NbWordsIndexesNoZero}},rb{{$i}}{{end}}); \
	SUBQ    q<>(SB), rb0; \
	{{- range $i := .NbWordsIndexesNoZero}}
	SBBQ  q<>+{{mul $i 8}}(SB), rb{{$i}}; \
	{{- end}}
	{{- range $i := .NbWordsIndexesFull}}
	CMOVQCC rb{{$i}}, ra{{$i}};  \
	MOVQ ra{{$i}}, {{mul $i 8}}(m0);  \
	{{- end}}

	

`

func (f *FFAmd64) GenerateDefines() {
	tmpl := template.Must(template.New("").
		Funcs(helpers()).
		Parse(tmplDefines))

	// execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, f); err != nil {
		panic(err)
	}

	f.WriteLn(buf.String())
}

func (f *FFAmd64) Mov(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {
	case []uint64:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1[i+o1], c2[i+o2])
			}
		}
	case amd64.Register:
		switch c2 := i2.(type) {
		case amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1.At(i+o1), c2.At(i+o2))
			}
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1.At(i+o1), c2[i+o2])
			}
		default:
			panic("unsupported")
		}
	case []amd64.Register:
		switch c2 := i2.(type) {
		case amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				f.MOVQ(c1[i+o1], c2.At(i+o2))
			}
		case []amd64.Register:
			f.copyElement(c1[o1:], c2[o2:])
			// for i := 0; i < f.NbWords; i++ {
			// 	f.MOVQ(c1[i+o1], c2[i+o2])
			// }
		default:
			panic("unsupported")
		}
	default:
		panic("unsupported")
	}

}

func (f *FFAmd64) Add(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {

	case amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.ADDQ(c1.At(i+o1), c2[i+o2])
				} else {
					f.ADCQ(c1.At(i+o1), c2[i+o2])
				}
			}
		}
	case []amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.ADDQ(c1[i+o1], c2[i+o2])
				} else {
					f.ADCQ(c1[i+o1], c2[i+o2])
				}
			}
		}
	default:
		panic("unsupported")
	}
}

func (f *FFAmd64) Sub(i1, i2 interface{}, offsets ...int) {
	var o1, o2 int
	if len(offsets) >= 1 {
		o1 = offsets[0]
		if len(offsets) >= 2 {
			o2 = offsets[1]
		}
	}
	switch c1 := i1.(type) {

	case amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.SUBQ(c1.At(i+o1), c2[i+o2])
				} else {
					f.SBBQ(c1.At(i+o1), c2[i+o2])
				}
			}
		}
	case []amd64.Register:
		switch c2 := i2.(type) {
		default:
			panic("unsupported")
		case []amd64.Register:
			for i := 0; i < f.NbWords; i++ {
				if i == 0 {
					f.SUBQ(c1[i+o1], c2[i+o2])
				} else {
					f.SBBQ(c1[i+o1], c2[i+o2])
				}
			}
		}
	default:
		panic("unsupported")
	}
}

func aggregate(values []string) string {
	var sb strings.Builder
	for _, v := range values {
		sb.WriteString(v)
	}
	return sb.String()
}

// Template helpers (txt/template)
func helpers() template.FuncMap {
	// functions used in template
	return template.FuncMap{
		"mul": mul,
		"imm": imm,
	}
}

func mul(a, b int) int {
	return a * b
}

func imm(t uint64) string {
	switch t {
	case 0:
		return "$0"
	case 1:
		return "$1"
	default:
		return fmt.Sprintf("$%#016x", t)
	}
}
