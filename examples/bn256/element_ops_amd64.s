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

#include "textflag.h"
#include "funcdata.h"

// modulus q
DATA q<>+0(SB)/8, $0x3c208c16d87cfd47
DATA q<>+8(SB)/8, $0x97816a916871ca8d
DATA q<>+16(SB)/8, $0xb85045b68181585d
DATA q<>+24(SB)/8, $0x30644e72e131a029
GLOBL q<>(SB), (RODATA+NOPTR), $32

// qInv0 q'[0]
DATA qInv0<>(SB)/8, $0x87d20782e4866389
GLOBL qInv0<>(SB), (RODATA+NOPTR), $8

#define REDUCE(ra0, ra1, ra2, ra3, rb0, rb1, rb2, rb3) \
	MOVQ    ra0, rb0;        \
	SUBQ    q<>(SB), ra0;    \
	MOVQ    ra1, rb1;        \
	SBBQ    q<>+8(SB), ra1;  \
	MOVQ    ra2, rb2;        \
	SBBQ    q<>+16(SB), ra2; \
	MOVQ    ra3, rb3;        \
	SBBQ    q<>+24(SB), ra3; \
	CMOVQCS rb0, ra0;        \
	CMOVQCS rb1, ra1;        \
	CMOVQCS rb2, ra2;        \
	CMOVQCS rb3, ra3;        \

// add(res, x, y *Element)
TEXT ·add(SB), NOSPLIT, $0-24
	MOVQ x+8(FP), AX
	MOVQ 0(AX), CX
	MOVQ 8(AX), BX
	MOVQ 16(AX), SI
	MOVQ 24(AX), DI
	MOVQ y+16(FP), DX
	ADDQ 0(DX), CX
	ADCQ 8(DX), BX
	ADCQ 16(DX), SI
	ADCQ 24(DX), DI

	// reduce element(CX,BX,SI,DI) using temp registers (R8,R9,R10,R11)
	REDUCE(CX,BX,SI,DI,R8,R9,R10,R11)

	MOVQ res+0(FP), R12
	MOVQ CX, 0(R12)
	MOVQ BX, 8(R12)
	MOVQ SI, 16(R12)
	MOVQ DI, 24(R12)
	RET

// sub(res, x, y *Element)
TEXT ·sub(SB), NOSPLIT, $0-24
	MOVQ    x+8(FP), SI
	MOVQ    0(SI), AX
	MOVQ    8(SI), DX
	MOVQ    16(SI), CX
	MOVQ    24(SI), BX
	MOVQ    y+16(FP), DI
	SUBQ    0(DI), AX
	SBBQ    8(DI), DX
	SBBQ    16(DI), CX
	SBBQ    24(DI), BX
	MOVQ    $0x3c208c16d87cfd47, R8
	MOVQ    $0x97816a916871ca8d, R9
	MOVQ    $0xb85045b68181585d, R10
	MOVQ    $0x30644e72e131a029, R11
	MOVQ    $0, R12
	CMOVQCC R12, R8
	CMOVQCC R12, R9
	CMOVQCC R12, R10
	CMOVQCC R12, R11
	ADDQ    R8, AX
	ADCQ    R9, DX
	ADCQ    R10, CX
	ADCQ    R11, BX
	MOVQ    res+0(FP), R13
	MOVQ    AX, 0(R13)
	MOVQ    DX, 8(R13)
	MOVQ    CX, 16(R13)
	MOVQ    BX, 24(R13)
	RET

// double(res, x *Element)
TEXT ·double(SB), NOSPLIT, $0-16
	MOVQ x+8(FP), AX
	MOVQ 0(AX), DX
	MOVQ 8(AX), CX
	MOVQ 16(AX), BX
	MOVQ 24(AX), SI
	ADDQ DX, DX
	ADCQ CX, CX
	ADCQ BX, BX
	ADCQ SI, SI

	// reduce element(DX,CX,BX,SI) using temp registers (DI,R8,R9,R10)
	REDUCE(DX,CX,BX,SI,DI,R8,R9,R10)

	MOVQ res+0(FP), R11
	MOVQ DX, 0(R11)
	MOVQ CX, 8(R11)
	MOVQ BX, 16(R11)
	MOVQ SI, 24(R11)
	RET

// neg(res, x *Element)
TEXT ·neg(SB), NOSPLIT, $0-16
	MOVQ  res+0(FP), DI
	MOVQ  x+8(FP), AX
	MOVQ  0(AX), DX
	MOVQ  8(AX), CX
	MOVQ  16(AX), BX
	MOVQ  24(AX), SI
	MOVQ  DX, AX
	ORQ   CX, AX
	ORQ   BX, AX
	ORQ   SI, AX
	TESTQ AX, AX
	JEQ   l1
	MOVQ  $0x3c208c16d87cfd47, R8
	SUBQ  DX, R8
	MOVQ  R8, 0(DI)
	MOVQ  $0x97816a916871ca8d, R8
	SBBQ  CX, R8
	MOVQ  R8, 8(DI)
	MOVQ  $0xb85045b68181585d, R8
	SBBQ  BX, R8
	MOVQ  R8, 16(DI)
	MOVQ  $0x30644e72e131a029, R8
	SBBQ  SI, R8
	MOVQ  R8, 24(DI)
	RET

l1:
	MOVQ AX, 0(DI)
	MOVQ AX, 8(DI)
	MOVQ AX, 16(DI)
	MOVQ AX, 24(DI)
	RET

// mul(res, x, y *Element)
TEXT ·mul(SB), $24-24

	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// however, to benefit from the ADCX and ADOX carry chains
	// we split the inner loops in 2:
	// for i=0 to N-1
	// 		for j=0 to N-1
	// 		    (A,t[j])  := t[j] + x[j]*y[i] + A
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 		    (C,t[j-1]) := t[j] + m*q[j] + C
	// 		t[N-1] = C + A

	NO_LOCAL_POINTERS
	CMPB ·supportAdx(SB), $1
	JNE  l2
	MOVQ x+8(FP), R14

	// x[0] = R15
	// x[1] = CX
	// x[2] = BX
	// x[3] = SI

	MOVQ 0(R14), R15
	MOVQ 8(R14), CX
	MOVQ 16(R14), BX
	MOVQ 24(R14), SI
	MOVQ y+16(FP), DI

	// t[0] = R8
	// t[1] = R9
	// t[2] = R10
	// t[3] = R11

	// clear the flags
	XORQ AX, AX
	MOVQ 0(DI), DX

	// (A,t[0])  := t[0] + x[0]*y[0] + A
	MULXQ R15, R8, R9

	// (A,t[1])  := t[1] + x[1]*y[0] + A
	MULXQ CX, AX, R10
	ADOXQ AX, R9

	// (A,t[2])  := t[2] + x[2]*y[0] + A
	MULXQ BX, AX, R11
	ADOXQ AX, R10

	// (A,t[3])  := t[3] + x[3]*y[0] + A
	MULXQ SI, AX, R12
	ADOXQ AX, R11

	// A += carries from ADCXQ and ADOXQ
	MOVQ  $0, AX
	ADOXQ AX, R12

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R8, DX

	// clear the flags
	XORQ AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R8, AX
	MOVQ  BP, R8

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R9, R8
	MULXQ q<>+8(SB), AX, R9
	ADOXQ AX, R8

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ R10, R9
	MULXQ q<>+16(SB), AX, R10
	ADOXQ AX, R9

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ R11, R10
	MULXQ q<>+24(SB), AX, R11
	ADOXQ AX, R10

	// t[3] = C + A
	MOVQ  $0, AX
	ADCXQ AX, R11
	ADOXQ R12, R11

	// clear the flags
	XORQ AX, AX
	MOVQ 8(DI), DX

	// (A,t[0])  := t[0] + x[0]*y[1] + A
	MULXQ R15, AX, R12
	ADOXQ AX, R8

	// (A,t[1])  := t[1] + x[1]*y[1] + A
	ADCXQ R12, R9
	MULXQ CX, AX, R12
	ADOXQ AX, R9

	// (A,t[2])  := t[2] + x[2]*y[1] + A
	ADCXQ R12, R10
	MULXQ BX, AX, R12
	ADOXQ AX, R10

	// (A,t[3])  := t[3] + x[3]*y[1] + A
	ADCXQ R12, R11
	MULXQ SI, AX, R12
	ADOXQ AX, R11

	// A += carries from ADCXQ and ADOXQ
	MOVQ  $0, AX
	ADCXQ AX, R12
	ADOXQ AX, R12

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R8, DX

	// clear the flags
	XORQ AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R8, AX
	MOVQ  BP, R8

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R9, R8
	MULXQ q<>+8(SB), AX, R9
	ADOXQ AX, R8

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ R10, R9
	MULXQ q<>+16(SB), AX, R10
	ADOXQ AX, R9

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ R11, R10
	MULXQ q<>+24(SB), AX, R11
	ADOXQ AX, R10

	// t[3] = C + A
	MOVQ  $0, AX
	ADCXQ AX, R11
	ADOXQ R12, R11

	// clear the flags
	XORQ AX, AX
	MOVQ 16(DI), DX

	// (A,t[0])  := t[0] + x[0]*y[2] + A
	MULXQ R15, AX, R12
	ADOXQ AX, R8

	// (A,t[1])  := t[1] + x[1]*y[2] + A
	ADCXQ R12, R9
	MULXQ CX, AX, R12
	ADOXQ AX, R9

	// (A,t[2])  := t[2] + x[2]*y[2] + A
	ADCXQ R12, R10
	MULXQ BX, AX, R12
	ADOXQ AX, R10

	// (A,t[3])  := t[3] + x[3]*y[2] + A
	ADCXQ R12, R11
	MULXQ SI, AX, R12
	ADOXQ AX, R11

	// A += carries from ADCXQ and ADOXQ
	MOVQ  $0, AX
	ADCXQ AX, R12
	ADOXQ AX, R12

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R8, DX

	// clear the flags
	XORQ AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R8, AX
	MOVQ  BP, R8

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R9, R8
	MULXQ q<>+8(SB), AX, R9
	ADOXQ AX, R8

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ R10, R9
	MULXQ q<>+16(SB), AX, R10
	ADOXQ AX, R9

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ R11, R10
	MULXQ q<>+24(SB), AX, R11
	ADOXQ AX, R10

	// t[3] = C + A
	MOVQ  $0, AX
	ADCXQ AX, R11
	ADOXQ R12, R11

	// clear the flags
	XORQ AX, AX
	MOVQ 24(DI), DX

	// (A,t[0])  := t[0] + x[0]*y[3] + A
	MULXQ R15, AX, R12
	ADOXQ AX, R8

	// (A,t[1])  := t[1] + x[1]*y[3] + A
	ADCXQ R12, R9
	MULXQ CX, AX, R12
	ADOXQ AX, R9

	// (A,t[2])  := t[2] + x[2]*y[3] + A
	ADCXQ R12, R10
	MULXQ BX, AX, R12
	ADOXQ AX, R10

	// (A,t[3])  := t[3] + x[3]*y[3] + A
	ADCXQ R12, R11
	MULXQ SI, AX, R12
	ADOXQ AX, R11

	// A += carries from ADCXQ and ADOXQ
	MOVQ  $0, AX
	ADCXQ AX, R12
	ADOXQ AX, R12

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R8, DX

	// clear the flags
	XORQ AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R8, AX
	MOVQ  BP, R8

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R9, R8
	MULXQ q<>+8(SB), AX, R9
	ADOXQ AX, R8

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ R10, R9
	MULXQ q<>+16(SB), AX, R10
	ADOXQ AX, R9

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ R11, R10
	MULXQ q<>+24(SB), AX, R11
	ADOXQ AX, R10

	// t[3] = C + A
	MOVQ  $0, AX
	ADCXQ AX, R11
	ADOXQ R12, R11

	// reduce element(R8,R9,R10,R11) using temp registers (R13,R14,R12,DI)
	REDUCE(R8,R9,R10,R11,R13,R14,R12,DI)

	MOVQ res+0(FP), AX
	MOVQ R8, 0(AX)
	MOVQ R9, 8(AX)
	MOVQ R10, 16(AX)
	MOVQ R11, 24(AX)
	RET

l2:
	MOVQ res+0(FP), AX
	MOVQ AX, (SP)
	MOVQ x+8(FP), AX
	MOVQ AX, 8(SP)
	MOVQ y+16(FP), AX
	MOVQ AX, 16(SP)
	CALL ·_mulGeneric(SB)
	RET

TEXT ·fromMont(SB), $8-8
	NO_LOCAL_POINTERS

	// the algorithm is described here
	// https://hackmd.io/@zkteam/modular_multiplication
	// when y = 1 we have:
	// for i=0 to N-1
	// 		t[i] = x[i]
	// for i=0 to N-1
	// 		m := t[0]*q'[0] mod W
	// 		C,_ := t[0] + m*q[0]
	// 		for j=1 to N-1
	// 		    (C,t[j-1]) := t[j] + m*q[j] + C
	// 		t[N-1] = C
	CMPB ·supportAdx(SB), $1
	JNE  l3
	MOVQ res+0(FP), DX
	MOVQ 0(DX), R14
	MOVQ 8(DX), R15
	MOVQ 16(DX), CX
	MOVQ 24(DX), BX
	XORQ DX, DX

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R14, DX
	XORQ  AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R14, AX
	MOVQ  BP, R14

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R15, R14
	MULXQ q<>+8(SB), AX, R15
	ADOXQ AX, R14

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ CX, R15
	MULXQ q<>+16(SB), AX, CX
	ADOXQ AX, R15

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ BX, CX
	MULXQ q<>+24(SB), AX, BX
	ADOXQ AX, CX
	MOVQ  $0, AX
	ADCXQ AX, BX
	ADOXQ AX, BX
	XORQ  DX, DX

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R14, DX
	XORQ  AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R14, AX
	MOVQ  BP, R14

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R15, R14
	MULXQ q<>+8(SB), AX, R15
	ADOXQ AX, R14

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ CX, R15
	MULXQ q<>+16(SB), AX, CX
	ADOXQ AX, R15

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ BX, CX
	MULXQ q<>+24(SB), AX, BX
	ADOXQ AX, CX
	MOVQ  $0, AX
	ADCXQ AX, BX
	ADOXQ AX, BX
	XORQ  DX, DX

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R14, DX
	XORQ  AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R14, AX
	MOVQ  BP, R14

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R15, R14
	MULXQ q<>+8(SB), AX, R15
	ADOXQ AX, R14

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ CX, R15
	MULXQ q<>+16(SB), AX, CX
	ADOXQ AX, R15

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ BX, CX
	MULXQ q<>+24(SB), AX, BX
	ADOXQ AX, CX
	MOVQ  $0, AX
	ADCXQ AX, BX
	ADOXQ AX, BX
	XORQ  DX, DX

	// m := t[0]*q'[0] mod W
	MOVQ  qInv0<>(SB), DX
	IMULQ R14, DX
	XORQ  AX, AX

	// C,_ := t[0] + m*q[0]
	MULXQ q<>+0(SB), AX, BP
	ADCXQ R14, AX
	MOVQ  BP, R14

	// (C,t[0]) := t[1] + m*q[1] + C
	ADCXQ R15, R14
	MULXQ q<>+8(SB), AX, R15
	ADOXQ AX, R14

	// (C,t[1]) := t[2] + m*q[2] + C
	ADCXQ CX, R15
	MULXQ q<>+16(SB), AX, CX
	ADOXQ AX, R15

	// (C,t[2]) := t[3] + m*q[3] + C
	ADCXQ BX, CX
	MULXQ q<>+24(SB), AX, BX
	ADOXQ AX, CX
	MOVQ  $0, AX
	ADCXQ AX, BX
	ADOXQ AX, BX

	// reduce element(R14,R15,CX,BX) using temp registers (SI,DI,R8,R9)
	REDUCE(R14,R15,CX,BX,SI,DI,R8,R9)

	MOVQ res+0(FP), AX
	MOVQ R14, 0(AX)
	MOVQ R15, 8(AX)
	MOVQ CX, 16(AX)
	MOVQ BX, 24(AX)
	RET

l3:
	MOVQ res+0(FP), AX
	MOVQ AX, (SP)
	CALL ·_fromMontGeneric(SB)
	RET

TEXT ·reduce(SB), NOSPLIT, $0-8
	MOVQ res+0(FP), AX
	MOVQ 0(AX), DX
	MOVQ 8(AX), CX
	MOVQ 16(AX), BX
	MOVQ 24(AX), SI

	// reduce element(DX,CX,BX,SI) using temp registers (DI,R8,R9,R10)
	REDUCE(DX,CX,BX,SI,DI,R8,R9,R10)

	MOVQ DX, 0(AX)
	MOVQ CX, 8(AX)
	MOVQ BX, 16(AX)
	MOVQ SI, 24(AX)
	RET

// MulBy3(x *Element)
TEXT ·MulBy3(SB), NOSPLIT, $0-8
	MOVQ x+0(FP), AX
	MOVQ 0(AX), DX
	MOVQ 8(AX), CX
	MOVQ 16(AX), BX
	MOVQ 24(AX), SI
	ADDQ DX, DX
	ADCQ CX, CX
	ADCQ BX, BX
	ADCQ SI, SI

	// reduce element(DX,CX,BX,SI) using temp registers (DI,R8,R9,R10)
	REDUCE(DX,CX,BX,SI,DI,R8,R9,R10)

	ADDQ 0(AX), DX
	ADCQ 8(AX), CX
	ADCQ 16(AX), BX
	ADCQ 24(AX), SI

	// reduce element(DX,CX,BX,SI) using temp registers (R11,R12,R13,R14)
	REDUCE(DX,CX,BX,SI,R11,R12,R13,R14)

	MOVQ DX, 0(AX)
	MOVQ CX, 8(AX)
	MOVQ BX, 16(AX)
	MOVQ SI, 24(AX)
	RET

// MulBy5(x *Element)
TEXT ·MulBy5(SB), NOSPLIT, $0-8
	MOVQ x+0(FP), AX
	MOVQ 0(AX), DX
	MOVQ 8(AX), CX
	MOVQ 16(AX), BX
	MOVQ 24(AX), SI
	ADDQ DX, DX
	ADCQ CX, CX
	ADCQ BX, BX
	ADCQ SI, SI

	// reduce element(DX,CX,BX,SI) using temp registers (DI,R8,R9,R10)
	REDUCE(DX,CX,BX,SI,DI,R8,R9,R10)

	ADDQ DX, DX
	ADCQ CX, CX
	ADCQ BX, BX
	ADCQ SI, SI

	// reduce element(DX,CX,BX,SI) using temp registers (R11,R12,R13,R14)
	REDUCE(DX,CX,BX,SI,R11,R12,R13,R14)

	ADDQ 0(AX), DX
	ADCQ 8(AX), CX
	ADCQ 16(AX), BX
	ADCQ 24(AX), SI

	// reduce element(DX,CX,BX,SI) using temp registers (R15,DI,R8,R9)
	REDUCE(DX,CX,BX,SI,R15,DI,R8,R9)

	MOVQ DX, 0(AX)
	MOVQ CX, 8(AX)
	MOVQ BX, 16(AX)
	MOVQ SI, 24(AX)
	RET
