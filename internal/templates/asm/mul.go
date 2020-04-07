package asm

const Mul = `

#include "textflag.h"

// func mulAsm{{.ElementName}}(res,y *{{.ElementName}})
// montgomery multiplication of res by y 
// stores the result in res
TEXT ·mulAsm{{.ElementName}}(SB), NOSPLIT, $0-16
	{{- /* do not change the order */ -}} 
	{{- $iReg := 0}}
	{{- $regt0 := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- $iReg = add 1 $iReg}}
	{{- end}}
	{{- $regX := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regY := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regA := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regC := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regM := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regYi := $iReg}}  {{- $iReg = add 1 $iReg}}

	// dereference our parameters
	MOVQ res+0(FP), {{reg $regX}}
	MOVQ y+8(FP), {{reg $regY}}

	// check if we support adx and mulx
	CMPB ·support_adx_{{.ElementName}}(SB), $1
	JNE no_adx
	 
	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// however, to benefit from the ADCX and ADOX carry chains
	// we split the inner loops in 2:
	// for i=0 to N-1
    // 		for j=0 to N-1
    // 		    (A,t[j])  := t[j] + a[j]*b[i] + A
    // 		m := t[0]*q'[0] mod W
    // 		C,_ := t[0] + m*q[0]
    // 		for j=1 to N-1
    // 		    (C,t[j-1]) := t[j] + m*q[j] + C
    // 		t[N-1] = C + A

	{{- range $i := .NbWordsIndexesFull}}
	// clear up the carry flags
	XORQ {{reg $regA}} , {{reg $regA}}

	// y[{{$i}}] in {{reg $regYi}}
	MOVQ {{mul $i 8}}({{reg $regY}}), {{reg $regYi}}

	// for j=0 to N-1
	//    (A,t[j])  := t[j] + x[j]*y[i] + A
	{{- range $j := $.NbWordsIndexesFull}}
		// res[{{$j}}] in DX
		MOVQ {{mul $j 8}}({{reg $regX}}), DX
		{{- if eq $j 0}}
			{{- if eq $i 0}}
				MULXQ {{reg $regYi}}, {{reg $regt0 $j}} ,  {{reg $regA}}
			{{- else}}
				MULXQ {{reg $regYi}}, AX,  {{reg $regA}}
				ADOXQ AX, {{reg $regt0 $j}} 
			{{- end}}
		{{- else}}
			{{- if eq $i 0}}
				MOVQ {{reg $regA}}, {{reg $regt0 $j}}
			{{- else}}
				ADCXQ {{reg $regA}}, {{reg $regt0 $j}}
			{{- end}}
			MULXQ {{reg $regYi}}, AX,  {{reg $regA}}
			ADOXQ AX, {{reg $regt0 $j}} 
		{{- end}}
	{{- end}}
	// add the last carries to {{reg $regA}} 
	MOVQ $0, DX
	ADCXQ DX, {{reg $regA}} 
	ADOXQ DX, {{reg $regA}} 
	
	// m := t[0]*q'[0] mod W
	MOVQ {{ $.ASMQInv0 }}, DX
	MULXQ {{reg $regt0}},{{reg $regM}}, DX

	// clear the carry flags
	XORQ DX, DX 

	// C,_ := t[0] + m*q[0]
	MOVQ {{ index $.ASMQ 0 }}, DX
	MULXQ {{reg $regM}}, AX, {{reg $regC}}
	ADCXQ {{reg $regt0}} ,AX

	// for j=1 to N-1
    //    (C,t[j-1]) := t[j] + m*q[j] + C
	{{- range $j := $.NbWordsIndexesNoZero}}
		{{- $k := sub $j 1}}
		
		MOVQ {{ index $.ASMQ $j }}, DX
		MULXQ {{reg $regM}}, AX, DX
		ADCXQ  {{reg $regt0 $j}}, {{reg $regC}} 
		ADOXQ AX, {{reg $regC}}
		MOVQ {{reg $regC}}, {{reg $regt0 $k}}
		{{- if eq $j $.NbWordsLastIndex}}
			MOVQ $0, AX
			ADCXQ AX, DX
			ADOXQ DX, {{reg $regA}}
			MOVQ {{reg $regA}}, {{reg $regt0 $.NbWordsLastIndex}}
		{{- else }}
			MOVQ DX, {{reg $regC}}
		{{- end}}
	{{- end}}

	{{- end}}

reduce:
	// reduce, constant time version
	// first we copy registers storing t in a separate set of registers
	// as SUBQ modifies the 2nd operand
	{{- /* registers after regY are not needed anymore */ -}}
	{{- /* u0 will be stored in DX */ -}}
	{{- $regu1 := $regY}}
	{{- $k := sub $.NbWords 1}}

	{{- /* temporary register to store moduli word for SBBQ */ -}}
	{{- $regQ := add $regY $k}}
	{{- range $i := .NbWordsIndexesFull}}
	{{- if eq $i 0}}
	MOVQ {{reg $regt0}}, DX
	{{- else}}
	{{- $k := sub $i 1}}
	MOVQ {{reg $regt0 $i}}, {{reg $regu1 $k}}
	{{- end}}
	{{- end }}

	{{- range $i := .NbWordsIndexesFull}}
	MOVQ {{ index $.ASMQ $i }}, {{reg $regQ}}
	{{- if eq $i 0}}
	SUBQ  {{reg $regQ}}, DX
	{{- else}}
	{{- $k := sub $i 1}}
	SBBQ  {{reg $regQ}}, {{reg $regu1 $k}}
	{{- end}}
	{{- end}}
	JCS t_is_smaller // no borrow, we return t
	MOVQ DX, ({{reg $regX}})
	{{- range $i := .NbWordsIndexesNoZero}}
	{{- $j := sub $i 1}}
	MOVQ {{reg $regu1 $j}}, {{mul $i 8}}({{reg $regX}})
	{{- end}}
	RET
t_is_smaller:
	{{- range $i := .NbWordsIndexesFull}}
	MOVQ {{reg $regt0 $i}}, {{mul $i 8}}({{reg $regX}})
	{{- end}}
	RET

no_adx:
	{{- range $i := .NbWordsIndexesFull}}
	// (A,t[0]) := t[0] + x[0]*y[{{$i}}]
	MOVQ ({{ reg $regX}}), AX // x[0]
	MOVQ {{mul $i 8}}({{reg $regY}}), {{reg $regYi}}
	MULQ {{reg $regYi}} // x[0] * y[{{$i}}]
	{{- if ne $i 0}}
	ADDQ AX, {{ reg $regt0}} 
	ADCQ $0, DX
	{{- end}}	
	MOVQ DX, {{ reg $regA}} 
	{{- if eq $i 0}}
	MOVQ AX, {{ reg $regt0}}
	{{end}}
	
	
	// m := t[0]*q'[0] mod W
	MOVQ {{ $.ASMQInv0 }}, {{ reg $regM}}
	IMULQ {{reg $regt0}} , {{ reg $regM}}

	// C,_ := t[0] + m*q[0]
	MOVQ {{ index $.ASMQ 0 }}, AX
	MULQ {{ reg $regM}}
	ADDQ {{ reg $regt0}} ,AX
	ADCQ $0, DX
	MOVQ  DX, {{ reg $regC}}

	// for j=1 to N-1
	//    (A,t[j])  := t[j] + x[j]*y[i] + A
    //    (C,t[j-1]) := t[j] + m*q[j] + C
	{{- range $j := $.NbWordsIndexesNoZero}}
	MOVQ {{mul $j 8}}({{ reg $regX}}), AX
	MULQ {{reg $regYi}} // x[{{$j}}] * y[{{$i}}]
	{{- if ne $i 0}}
	ADDQ {{ reg $regA}}, {{reg $regt0 $j}}
	ADCQ $0, DX
	ADDQ AX, {{reg $regt0 $j}}
	ADCQ $0, DX
	{{- else}}
	MOVQ {{ reg $regA}}, {{reg $regt0 $j}}
	ADDQ AX, {{reg $regt0 $j}}
	ADCQ $0, DX
	{{- end}}
	MOVQ DX, {{ reg $regA}}

	MOVQ {{ index $.ASMQ $j }}, AX
	MULQ {{ reg $regM}}
	ADDQ  {{reg $regt0 $j}}, {{ reg $regC}}
	ADCQ $0, DX
	ADDQ AX, {{ reg $regC}}
	ADCQ $0, DX
	{{$k := sub $j 1}}
	MOVQ {{ reg $regC}}, {{reg $regt0 $k}}
	MOVQ DX, {{ reg $regC}}
	{{- end}}

	ADDQ {{ reg $regC}}, {{ reg $regA}}
	MOVQ {{ reg $regA}}, {{reg $regt0 $.NbWordsLastIndex}}

	{{- end}}

	JMP reduce
`
