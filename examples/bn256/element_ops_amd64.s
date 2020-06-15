#include "textflag.h"
TEXT ·mulElement(SB), NOSPLIT, $0-24

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
	
    MOVQ x+8(FP), DI                                       // dereference x
    CMPB ·supportAdx(SB), $0x0000000000000001             // check if we support MULX and ADOX instructions
    JNE no_adx                                            // no support for MULX or ADOX instructions
    MOVQ y+16(FP), R10                                     // dereference y
    MOVQ 0(DI), R11                                        // R11 = x[0]
    MOVQ 8(DI), R12                                        // R12 = x[1]
    MOVQ 16(DI), R13                                       // R13 = x[2]
    MOVQ 24(DI), R14                                       // R14 = x[3]
    // outter loop 0
    XORQ DX, DX                                            // clear up flags
    MOVQ 0(R10), DX                                        // DX = y[0]
    MULXQ R11, CX, BX                                       // t[0], t[1] = y[0] * x[0]
    MULXQ R12, AX, BP
    ADOXQ AX, BX
    MULXQ R13, AX, SI
    ADOXQ AX, BP
    MULXQ R14, AX, R9
    ADOXQ AX, SI
    // add the last carries to R9
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    // outter loop 1
    XORQ DX, DX                                            // clear up flags
    MOVQ 8(R10), DX                                        // DX = y[1]
    MULXQ R11, AX, R9
    ADOXQ AX, CX
    ADCXQ R9, BX                                            // t[1] += regA
    MULXQ R12, AX, R9
    ADOXQ AX, BX
    ADCXQ R9, BP                                            // t[2] += regA
    MULXQ R13, AX, R9
    ADOXQ AX, BP
    ADCXQ R9, SI                                            // t[3] += regA
    MULXQ R14, AX, R9
    ADOXQ AX, SI
    // add the last carries to R9
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    // outter loop 2
    XORQ DX, DX                                            // clear up flags
    MOVQ 16(R10), DX                                       // DX = y[2]
    MULXQ R11, AX, R9
    ADOXQ AX, CX
    ADCXQ R9, BX                                            // t[1] += regA
    MULXQ R12, AX, R9
    ADOXQ AX, BX
    ADCXQ R9, BP                                            // t[2] += regA
    MULXQ R13, AX, R9
    ADOXQ AX, BP
    ADCXQ R9, SI                                            // t[3] += regA
    MULXQ R14, AX, R9
    ADOXQ AX, SI
    // add the last carries to R9
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    // outter loop 3
    XORQ DX, DX                                            // clear up flags
    MOVQ 24(R10), DX                                       // DX = y[3]
    MULXQ R11, AX, R9
    ADOXQ AX, CX
    ADCXQ R9, BX                                            // t[1] += regA
    MULXQ R12, AX, R9
    ADOXQ AX, BX
    ADCXQ R9, BP                                            // t[2] += regA
    MULXQ R13, AX, R9
    ADOXQ AX, BP
    ADCXQ R9, SI                                            // t[3] += regA
    MULXQ R14, AX, R9
    ADOXQ AX, SI
    // add the last carries to R9
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R9
    ADOXQ DX, R9
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R9, SI
    MOVQ res+0(FP), DI                                     // dereference res
    MOVQ CX, R15
    SUBQ ·qElement+0(SB), R15
    MOVQ BX, R8
    SBBQ ·qElement+8(SB), R8
    MOVQ BP, R10
    SBBQ ·qElement+16(SB), R10
    MOVQ SI, R9
    SBBQ ·qElement+24(SB), R9
    CMOVQCC R15, CX
    CMOVQCC R8, BX
    CMOVQCC R10, BP
    CMOVQCC R9, SI
    MOVQ CX, 0(DI)
    MOVQ BX, 8(DI)
    MOVQ BP, 16(DI)
    MOVQ SI, 24(DI)
    RET
no_adx:
    MOVQ y+16(FP), R15                                     // dereference y
    MOVQ 0(DI), AX
    MOVQ 0(R15), R12
    MULQ R12
    MOVQ AX, CX
    MOVQ DX, R13
    MOVQ $0x87d20782e4866389, R14
    IMULQ CX, R14
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ 8(DI), AX
    MULQ R12
    MOVQ R13, BX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12
    MOVQ R13, BP
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12
    MOVQ R13, SI
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    ADDQ R11, R13
    MOVQ R13, SI
    MOVQ 0(DI), AX
    MOVQ 8(R15), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x87d20782e4866389, R14
    IMULQ CX, R14
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ 8(DI), AX
    MULQ R12
    ADDQ R13, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12
    ADDQ R13, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12
    ADDQ R13, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    ADDQ R11, R13
    MOVQ R13, SI
    MOVQ 0(DI), AX
    MOVQ 16(R15), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x87d20782e4866389, R14
    IMULQ CX, R14
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ 8(DI), AX
    MULQ R12
    ADDQ R13, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12
    ADDQ R13, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12
    ADDQ R13, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    ADDQ R11, R13
    MOVQ R13, SI
    MOVQ 0(DI), AX
    MOVQ 24(R15), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x87d20782e4866389, R14
    IMULQ CX, R14
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ 8(DI), AX
    MULQ R12
    ADDQ R13, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ 16(DI), AX
    MULQ R12
    ADDQ R13, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ 24(DI), AX
    MULQ R12
    ADDQ R13, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    ADDQ R11, R13
    MOVQ R13, SI
    MOVQ res+0(FP), DI                                     // dereference res
    MOVQ CX, R8
    SUBQ ·qElement+0(SB), R8
    MOVQ BX, R10
    SBBQ ·qElement+8(SB), R10
    MOVQ BP, R9
    SBBQ ·qElement+16(SB), R9
    MOVQ SI, R11
    SBBQ ·qElement+24(SB), R11
    CMOVQCC R8, CX
    CMOVQCC R10, BX
    CMOVQCC R9, BP
    CMOVQCC R11, SI
    MOVQ CX, 0(DI)
    MOVQ BX, 8(DI)
    MOVQ BP, 16(DI)
    MOVQ SI, 24(DI)
    RET

TEXT ·fromMontElement(SB), NOSPLIT, $0-8

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
    MOVQ res+0(FP), DI                                     // dereference x
    MOVQ 0(DI), CX                                         // t[0] = x[0]
    MOVQ 8(DI), BX                                         // t[1] = x[1]
    MOVQ 16(DI), BP                                        // t[2] = x[2]
    MOVQ 24(DI), SI                                        // t[3] = x[3]
    CMPB ·supportAdx(SB), $0x0000000000000001             // check if we support MULX and ADOX instructions
    JNE no_adx                                            // no support for MULX or ADOX instructions
    // outter loop 0
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ AX, SI
    // outter loop 1
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ AX, SI
    // outter loop 2
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ AX, SI
    // outter loop 3
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·qElement+0(SB), AX, R8
    ADCXQ CX, AX
    MOVQ R8, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ AX, SI
    MOVQ CX, R9
    SUBQ ·qElement+0(SB), R9
    MOVQ BX, R10
    SBBQ ·qElement+8(SB), R10
    MOVQ BP, R11
    SBBQ ·qElement+16(SB), R11
    MOVQ SI, R12
    SBBQ ·qElement+24(SB), R12
    CMOVQCC R9, CX
    CMOVQCC R10, BX
    CMOVQCC R11, BP
    CMOVQCC R12, SI
    MOVQ CX, 0(DI)
    MOVQ BX, 8(DI)
    MOVQ BP, 16(DI)
    MOVQ SI, 24(DI)
    RET
no_adx:
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ R13, SI
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ R13, SI
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ R13, SI
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ R13, SI
    MOVQ CX, R10
    SUBQ ·qElement+0(SB), R10
    MOVQ BX, R11
    SBBQ ·qElement+8(SB), R11
    MOVQ BP, R12
    SBBQ ·qElement+16(SB), R12
    MOVQ SI, R13
    SBBQ ·qElement+24(SB), R13
    CMOVQCC R10, CX
    CMOVQCC R11, BX
    CMOVQCC R12, BP
    CMOVQCC R13, SI
    MOVQ CX, 0(DI)
    MOVQ BX, 8(DI)
    MOVQ BP, 16(DI)
    MOVQ SI, 24(DI)
    RET

TEXT ·squareElement(SB), NOSPLIT, $0-16

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

	// if adx and mulx instructions are not available, uses MUL algorithm.
	
    CMPB ·supportAdx(SB), $0x0000000000000001             // check if we support MULX and ADOX instructions
    JNE no_adx                                            // no support for MULX or ADOX instructions
    MOVQ y+8(FP), DI                                       // dereference y
    // outter loop 0
    XORQ AX, AX                                            // clear up flags
    // dx = y[0]
    MOVQ 0(DI), DX
    MULXQ 8(DI), R9, R10
    MULXQ 16(DI), AX, R11
    ADCXQ AX, R10
    MULXQ 24(DI), AX, R8
    ADCXQ AX, R11
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, CX, DX
    ADCXQ R9, R9
    MOVQ R9, BX
    ADOXQ DX, BX
    ADCXQ R10, R10
    MOVQ R10, BP
    ADOXQ AX, BP
    ADCXQ R11, R11
    MOVQ R11, SI
    ADOXQ AX, SI
    ADCXQ R8, R8
    ADOXQ AX, R8
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·qElement+0(SB), AX, R12
    ADCXQ CX, AX
    MOVQ R12, CX
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R8, SI
    // outter loop 1
    XORQ AX, AX                                            // clear up flags
    // dx = y[1]
    MOVQ 8(DI), DX
    MULXQ 16(DI), R13, R14
    MULXQ 24(DI), AX, R8
    ADCXQ AX, R14
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    XORQ AX, AX                                            // clear up flags
    ADCXQ R13, R13
    ADOXQ R13, BP
    ADCXQ R14, R14
    ADOXQ R14, SI
    ADCXQ R8, R8
    ADOXQ AX, R8
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, AX, DX
    ADOXQ AX, BX
    MOVQ $0x0000000000000000, AX
    ADOXQ DX, BP
    ADOXQ AX, SI
    ADOXQ AX, R8
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·qElement+0(SB), AX, R15
    ADCXQ CX, AX
    MOVQ R15, CX
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R8, SI
    // outter loop 2
    XORQ AX, AX                                            // clear up flags
    // dx = y[2]
    MOVQ 16(DI), DX
    MULXQ 24(DI), R9, R8
    ADCXQ R9, R9
    ADOXQ R9, SI
    ADCXQ R8, R8
    ADOXQ AX, R8
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, AX, DX
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADOXQ DX, SI
    ADOXQ AX, R8
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·qElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R8, SI
    // outter loop 3
    XORQ AX, AX                                            // clear up flags
    // dx = y[3]
    MOVQ 24(DI), DX
    MULXQ DX, AX, R8
    ADCXQ AX, SI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    MOVQ CX, DX
    MULXQ ·qElementInv0(SB), DX, AX                        // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·qElement+0(SB), AX, R11
    ADCXQ CX, AX
    MOVQ R11, CX
    ADCXQ BX, CX
    MULXQ ·qElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·qElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·qElement+24(SB), AX, SI
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, SI
    ADOXQ R8, SI
    // dereference res
    MOVQ res+0(FP), R12
    MOVQ CX, R13
    SUBQ ·qElement+0(SB), R13
    MOVQ BX, R14
    SBBQ ·qElement+8(SB), R14
    MOVQ BP, R15
    SBBQ ·qElement+16(SB), R15
    MOVQ SI, R9
    SBBQ ·qElement+24(SB), R9
    CMOVQCC R13, CX
    CMOVQCC R14, BX
    CMOVQCC R15, BP
    CMOVQCC R9, SI
    MOVQ CX, 0(R12)
    MOVQ BX, 8(R12)
    MOVQ BP, 16(R12)
    MOVQ SI, 24(R12)
    RET
no_adx:
    // dereference y
    MOVQ y+8(FP), R13
    MOVQ 0(R13), AX
    MOVQ 0(R13), R11
    MULQ R11
    MOVQ AX, CX
    MOVQ DX, DI
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R13), AX
    MULQ R11
    MOVQ DI, BX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R13), AX
    MULQ R11
    MOVQ DI, BP
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R13), AX
    MULQ R11
    MOVQ DI, SI
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    ADDQ R10, DI
    MOVQ DI, SI
    MOVQ 0(R13), AX
    MOVQ 8(R13), R11
    MULQ R11
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R13), AX
    MULQ R11
    ADDQ DI, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R13), AX
    MULQ R11
    ADDQ DI, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R13), AX
    MULQ R11
    ADDQ DI, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    ADDQ R10, DI
    MOVQ DI, SI
    MOVQ 0(R13), AX
    MOVQ 16(R13), R11
    MULQ R11
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R13), AX
    MULQ R11
    ADDQ DI, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R13), AX
    MULQ R11
    ADDQ DI, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R13), AX
    MULQ R11
    ADDQ DI, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    ADDQ R10, DI
    MOVQ DI, SI
    MOVQ 0(R13), AX
    MOVQ 24(R13), R11
    MULQ R11
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x87d20782e4866389, R8
    IMULQ CX, R8
    MOVQ $0x3c208c16d87cfd47, AX
    MULQ R8
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R13), AX
    MULQ R11
    ADDQ DI, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x97816a916871ca8d, AX
    MULQ R8
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R13), AX
    MULQ R11
    ADDQ DI, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0xb85045b68181585d, AX
    MULQ R8
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R13), AX
    MULQ R11
    ADDQ DI, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, DI
    MOVQ $0x30644e72e131a029, AX
    MULQ R8
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    ADDQ R10, DI
    MOVQ DI, SI
    // dereference res
    MOVQ res+0(FP), R12
    MOVQ CX, R14
    SUBQ ·qElement+0(SB), R14
    MOVQ BX, R15
    SBBQ ·qElement+8(SB), R15
    MOVQ BP, R9
    SBBQ ·qElement+16(SB), R9
    MOVQ SI, R10
    SBBQ ·qElement+24(SB), R10
    CMOVQCC R14, CX
    CMOVQCC R15, BX
    CMOVQCC R9, BP
    CMOVQCC R10, SI
    MOVQ CX, 0(R12)
    MOVQ BX, 8(R12)
    MOVQ BP, 16(R12)
    MOVQ SI, 24(R12)
    RET

TEXT ·reduceElement(SB), NOSPLIT, $0-8
    MOVQ res+0(FP), CX                                     // dereference x
    MOVQ 0(CX), BX                                         // t[0] = x[0]
    MOVQ 8(CX), BP                                         // t[1] = x[1]
    MOVQ 16(CX), SI                                        // t[2] = x[2]
    MOVQ 24(CX), DI                                        // t[3] = x[3]
    MOVQ BX, R8
    SUBQ ·qElement+0(SB), R8
    MOVQ BP, R9
    SBBQ ·qElement+8(SB), R9
    MOVQ SI, R10
    SBBQ ·qElement+16(SB), R10
    MOVQ DI, R11
    SBBQ ·qElement+24(SB), R11
    CMOVQCC R8, BX
    CMOVQCC R9, BP
    CMOVQCC R10, SI
    CMOVQCC R11, DI
    MOVQ BX, 0(CX)
    MOVQ BP, 8(CX)
    MOVQ SI, 16(CX)
    MOVQ DI, 24(CX)
    RET

TEXT ·addElement(SB), NOSPLIT, $0-24
    MOVQ x+8(FP), CX                                       // dereference x
    MOVQ y+16(FP), BX                                      // dereference y
    MOVQ 0(CX), BP                                         // t[0] = x[0]
    MOVQ 8(CX), SI                                         // t[1] = x[1]
    MOVQ 16(CX), DI                                        // t[2] = x[2]
    MOVQ 24(CX), R8                                        // t[3] = x[3]
    ADDQ 0(BX), BP
    ADCQ 8(BX), SI
    ADCQ 16(BX), DI
    ADCQ 24(BX), R8
    // note that we don't check for the carry here, as this code was generated assuming F.NoCarry condition is set (see goff for more details)
    MOVQ res+0(FP), CX                                     // dereference res
    MOVQ BP, R9
    SUBQ ·qElement+0(SB), R9
    MOVQ SI, R10
    SBBQ ·qElement+8(SB), R10
    MOVQ DI, R11
    SBBQ ·qElement+16(SB), R11
    MOVQ R8, R12
    SBBQ ·qElement+24(SB), R12
    CMOVQCC R9, BP
    CMOVQCC R10, SI
    CMOVQCC R11, DI
    CMOVQCC R12, R8
    MOVQ BP, 0(CX)
    MOVQ SI, 8(CX)
    MOVQ DI, 16(CX)
    MOVQ R8, 24(CX)
    RET

TEXT ·doubleElement(SB), NOSPLIT, $0-16
    MOVQ res+0(FP), CX                                     // dereference x
    MOVQ y+8(FP), BX                                       // dereference y
    MOVQ 0(BX), BP                                         // t[0] = y[0]
    MOVQ 8(BX), SI                                         // t[1] = y[1]
    MOVQ 16(BX), DI                                        // t[2] = y[2]
    MOVQ 24(BX), R8                                        // t[3] = y[3]
    ADDQ BP, BP
    ADCQ SI, SI
    ADCQ DI, DI
    ADCQ R8, R8
    // note that we don't check for the carry here, as this code was generated assuming F.NoCarry condition is set (see goff for more details)
    MOVQ BP, R9
    SUBQ ·qElement+0(SB), R9
    MOVQ SI, R10
    SBBQ ·qElement+8(SB), R10
    MOVQ DI, R11
    SBBQ ·qElement+16(SB), R11
    MOVQ R8, R12
    SBBQ ·qElement+24(SB), R12
    CMOVQCC R9, BP
    CMOVQCC R10, SI
    CMOVQCC R11, DI
    CMOVQCC R12, R8
    MOVQ BP, 0(CX)
    MOVQ SI, 8(CX)
    MOVQ DI, 16(CX)
    MOVQ R8, 24(CX)
    RET

TEXT ·subElement(SB), NOSPLIT, $0-24
    MOVQ x+8(FP), DI                                       // dereference x
    MOVQ y+16(FP), R8                                      // dereference y
    MOVQ 0(DI), CX                                         // t[0] = x[0]
    MOVQ 8(DI), BX                                         // t[1] = x[1]
    MOVQ 16(DI), BP                                        // t[2] = x[2]
    MOVQ 24(DI), SI                                        // t[3] = x[3]
    XORQ DX, DX
    SUBQ 0(R8), CX
    SBBQ 8(R8), BX
    SBBQ 16(R8), BP
    SBBQ 24(R8), SI
    MOVQ $0x3c208c16d87cfd47, R9
    MOVQ $0x97816a916871ca8d, R10
    MOVQ $0xb85045b68181585d, R11
    MOVQ $0x30644e72e131a029, R12
    CMOVQCC DX, R9
    CMOVQCC DX, R10
    CMOVQCC DX, R11
    CMOVQCC DX, R12
    ADDQ R9, CX
    ADCQ R10, BX
    ADCQ R11, BP
    ADCQ R12, SI
    MOVQ res+0(FP), DI                                     // dereference res
    MOVQ CX, 0(DI)
    MOVQ BX, 8(DI)
    MOVQ BP, 16(DI)
    MOVQ SI, 24(DI)
    RET
