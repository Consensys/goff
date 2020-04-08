package asm

// TODO need to add the no_adx version in unit tests
const Square = `

#include "textflag.h"

// func Square{{.ElementName}}(res,x *{{.ElementName}})
// montgomery squaring of x
// stores the result in res
TEXT ·Square{{.ElementName}}(SB), NOSPLIT, $0-16

	{{- /* do not change the order */ -}} 
	{{- $iReg := 0}}
	{{- $regt0 := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- $iReg = add 1 $iReg}}
	{{- end}}
	{{- $regX := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regA := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regC := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regM := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regXi := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regP := $iReg}}  {{- $iReg = add 1 $iReg}}
	{{- $regSuperHi := $iReg}}  {{- $iReg = add 1 $iReg}}

	// dereference x
	MOVQ x+8(FP), {{reg $regX}}



	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// for i=0 to N-1
    // A, t[i] = x[i] * x[i] + t[i]
    // p = 0
    // for j=i+1 to N-1
    //     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
    // m = t[0] * q'[0]
    // C, _ = t[0] + q[0]*m
    // for j=1 to N-1
    //     C, t[j-1] = q[j]*m +  t[j] + C
	// t[N-1] = C + A
	
	// check if we support adx and mulx
	// CMPB ·supportAdx(SB), $1
	// JNE no_adx

	{{- range $i := .NbWordsIndexesFull}}
		// XORQ {{reg $regt0 $i}}, {{reg $regt0 $i}}
	{{- end}}

	// for i=0 to N-1
	{{- range $i := .NbWordsIndexesFull}}

		// ---------------------------------------------------------------------------------------------
		// outter loop {{$i}}

		// clear up the carry flags
		XORQ {{reg $regP}}, {{reg $regP}}

		// A, t[{{$i}}] = x[{{$i}}] * x[{{$i}}] + t[{{$i}}]
		MOVQ {{mul $i 8}}({{reg $regX}}), {{reg $regXi}} 
		MOVQ {{reg $regXi}}, DX
		MULXQ DX, AX, {{reg $regA}}   // x[{{$i}}] * x[{{$i}}]
		{{if ne $i 0}}
			ADCXQ AX, {{reg $regt0 $i}} 
		{{else}}
			MOVQ AX, {{reg $regt0 $i}} 
		{{end}}
		
		// for j=i+1 to N-1
		//     A,t[j] = x[j]*x[i] + t[j] + A
		{{- range $j := $.NbWordsIndexesNoZero}}
			{{- $doubling_round := gt $j $i}}
			{{- if $doubling_round}}
				MOVQ {{mul $j 8}}({{reg $regX}}), DX
				{{- if eq $i 0}}
					MOVQ {{reg $regA}}, {{reg $regt0 $j}}
				{{- else}}
					ADCXQ {{reg $regA}}, {{reg $regt0 $j}}
				{{- end}}
				MULXQ {{reg $regXi}}, AX,  {{reg $regA}}
				ADOXQ AX, {{reg $regt0 $j}} 
			{{- end}}
		{{- end}}
		// add the last carries to {{reg $regA}} 
		MOVQ $0, DX
		ADCXQ DX, {{reg $regA}} 
		ADOXQ DX, {{reg $regA}} 

		XORQ {{reg $regC}}, {{reg $regC}}

		{{- range $j := $.NbWordsIndexesNoZero}}
			{{- $doubling_round := gt $j $i}}
			{{- if $doubling_round}}
				MOVQ {{mul $j 8}}({{reg $regX}}), DX
				ADCXQ {{reg $regC}}, {{reg $regt0 $j}}
				MULXQ {{reg $regXi}}, AX,  {{reg $regC}}
				ADOXQ AX, {{reg $regt0 $j}} 
			{{- end}}
		{{- end}}

		MOVQ $0, DX
		ADOXQ DX, {{reg $regC}} 
		ADCXQ {{reg $regC}}, {{reg $regA}}


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
	// dereference result
	MOVQ res+0(FP), AX
	// reduce, constant time version
	// first we copy registers storing t in a separate set of registers
	// as SUBQ modifies the 2nd operand
	{{- /* registers after regX are not needed anymore */ -}}
	{{- /* u0 will be stored in DX */ -}}
	{{- $regu1 := $regX}}
	{{- $k := sub $.NbWords 1}}

	{{- /* temporary register to store moduli word for SBBQ */ -}}
	{{- $regQ := add $regX $k}}
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

	// borrow is set, we return u
	MOVQ DX, (AX)
	{{- range $i := .NbWordsIndexesNoZero}}
		{{- $j := sub $i 1}}
		MOVQ {{reg $regu1 $j}}, {{mul $i 8}}(AX)
	{{- end}}
	RET

t_is_smaller:
	{{- range $i := .NbWordsIndexesFull}}
		MOVQ {{reg $regt0 $i}}, {{mul $i 8}}(AX)
	{{- end}}
	RET

no_adx:
	// for i=0 to N-1
	{{- range $i := .NbWordsIndexesFull}}

		// ---------------------------------------------------------------------------------------------
		// outter loop {{$i}}

		// A, t[{{$i}}] = x[{{$i}}] * x[{{$i}}] + t[{{$i}}]
		MOVQ {{mul $i 8}}({{reg $regX}}), {{reg $regXi}} 
		MOVQ {{reg $regXi}}, AX
		MULQ AX // x[{{$i}}] * x[{{$i}}]
		{{if ne $i 0}}
			ADDQ AX, {{reg $regt0 $i}} 
			ADCQ $0, DX
		{{else}}
			MOVQ AX, {{reg $regt0 $i}} 
		{{end}}
		MOVQ DX, {{reg $regA}} 
		XORQ {{reg $regP}}, {{reg $regP}}
		

		// for j=i+1 to N-1
		//     p,A,t[j] = 2*x[j]*x[i] + t[j] + (p,A)
		{{- range $j := $.NbWordsIndexesNoZero}}
			{{- $doubling_round := gt $j $i}}
			{{- if $doubling_round}}
				XORQ {{reg $regSuperHi}}, {{reg $regSuperHi}}
				MOVQ {{mul $j 8}}({{reg $regX}}), AX
				MULQ {{reg $regXi}}
				ADDQ AX, AX
				ADCQ DX, DX
				ADCQ $0, {{reg $regSuperHi}}

				{{- if ne $i 0}}
					ADDQ {{reg $regt0 $j}}, {{reg $regA}}
					ADCQ $0, DX
				{{- end}}
				
				ADDQ {{reg $regA}}, AX
				ADCQ {{reg $regP}}, DX
				
				MOVQ {{reg $regSuperHi}}, {{reg $regP}}
				MOVQ DX, {{reg $regA}}
				MOVQ AX, {{reg $regt0 $j}}
			{{- end}}
		{{- end}}

		// m = t[0] * q'[0]
		MOVQ {{ $.ASMQInv0 }}, {{reg $regM}}
		IMULQ {{reg $regt0}} , {{reg $regM}}

		// C, _ = t[0] + q[0]*m
		MOVQ {{ index $.ASMQ 0 }}, AX
		MULQ {{reg $regM}}
		ADDQ {{reg $regt0}} ,AX
		ADCQ $0, DX
		MOVQ  DX, {{reg $regC}}

		// for j=1 to N-1
		//     C, t[j-1] = q[j]*m +  t[j] + C
		{{- range $j := $.NbWordsIndexesNoZero}}
			MOVQ {{ index $.ASMQ $j }}, AX
			MULQ {{reg $regM}}
			ADDQ  {{reg $regt0 $j}}, {{reg $regC}}
			ADCQ $0, DX
			ADDQ AX, {{reg $regC}}
			ADCQ $0, DX
			{{$k := sub $j 1}}
			MOVQ {{reg $regC}}, {{reg $regt0 $k}}
			MOVQ DX, {{reg $regC}}
		{{- end}}

		// t[N-1] = C + A
		ADDQ {{reg $regC}}, {{reg $regA}}
		MOVQ {{reg $regA}}, {{reg $regt0 $.NbWordsLastIndex}}
	{{- end}}

	JMP reduce
`
