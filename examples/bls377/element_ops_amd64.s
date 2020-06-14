#include "textflag.h"
TEXT ·mulAssignElement(SB), NOSPLIT, $0-16

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
	
    MOVQ res+0(FP), R9                                     // dereference x
    CMPB ·supportAdx(SB), $0x0000000000000001             // check if we support MULX and ADOX instructions
    JNE no_adx                                            // no support for MULX or ADOX instructions
    MOVQ y+8(FP), R12                                      // dereference y
    MOVQ 0(R9), R13                                        // R13 = x[0]
    MOVQ 8(R9), R14                                        // R14 = x[1]
    MOVQ 16(R9), R15                                       // R15 = x[2]
    // outter loop 0
    XORQ DX, DX                                            // clear up flags
    MOVQ 0(R12), DX                                        // DX = y[0]
    MULXQ R13, CX, BX                                       // t[0], t[1] = y[0] * x[0]
    MULXQ R14, AX, BP
    ADOXQ AX, BX
    MULXQ R15, AX, SI
    ADOXQ AX, BP
    MULXQ 24(R9), AX, DI
    ADOXQ AX, SI
    MULXQ 32(R9), AX, R8
    ADOXQ AX, DI
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 1
    XORQ DX, DX                                            // clear up flags
    MOVQ 8(R12), DX                                        // DX = y[1]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 2
    XORQ DX, DX                                            // clear up flags
    MOVQ 16(R12), DX                                       // DX = y[2]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 3
    XORQ DX, DX                                            // clear up flags
    MOVQ 24(R12), DX                                       // DX = y[3]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 4
    XORQ DX, DX                                            // clear up flags
    MOVQ 32(R12), DX                                       // DX = y[4]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 5
    XORQ DX, DX                                            // clear up flags
    MOVQ 40(R12), DX                                       // DX = y[5]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    MOVQ CX, R10
    SUBQ ·modulusElement+0(SB), R10
    MOVQ BX, R12
    SBBQ ·modulusElement+8(SB), R12
    MOVQ BP, R11
    SBBQ ·modulusElement+16(SB), R11
    MOVQ SI, R13
    SBBQ ·modulusElement+24(SB), R13
    MOVQ DI, R14
    SBBQ ·modulusElement+32(SB), R14
    MOVQ R8, R15
    SBBQ ·modulusElement+40(SB), R15
    CMOVQCC R10, CX
    CMOVQCC R12, BX
    CMOVQCC R11, BP
    CMOVQCC R13, SI
    CMOVQCC R14, DI
    CMOVQCC R15, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET
no_adx:
    MOVQ y+8(FP), R14                                      // dereference y
    MOVQ 0(R9), AX
    MOVQ 0(R14), R12
    MULQ R12
    MOVQ AX, CX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    MOVQ R11, BX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    MOVQ R11, BP
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    MOVQ R11, SI
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    MOVQ R11, DI
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    MOVQ R11, R8
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 8(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 16(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 24(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 32(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 40(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ CX, R15
    SUBQ ·modulusElement+0(SB), R15
    MOVQ BX, R10
    SBBQ ·modulusElement+8(SB), R10
    MOVQ BP, R12
    SBBQ ·modulusElement+16(SB), R12
    MOVQ SI, R11
    SBBQ ·modulusElement+24(SB), R11
    MOVQ DI, R13
    SBBQ ·modulusElement+32(SB), R13
    MOVQ R8, R14
    SBBQ ·modulusElement+40(SB), R14
    CMOVQCC R15, CX
    CMOVQCC R10, BX
    CMOVQCC R12, BP
    CMOVQCC R11, SI
    CMOVQCC R13, DI
    CMOVQCC R14, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET

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
	
    MOVQ x+8(FP), R9                                       // dereference x
    CMPB ·supportAdx(SB), $0x0000000000000001             // check if we support MULX and ADOX instructions
    JNE no_adx                                            // no support for MULX or ADOX instructions
    MOVQ y+16(FP), R12                                     // dereference y
    MOVQ 0(R9), R13                                        // R13 = x[0]
    MOVQ 8(R9), R14                                        // R14 = x[1]
    MOVQ 16(R9), R15                                       // R15 = x[2]
    // outter loop 0
    XORQ DX, DX                                            // clear up flags
    MOVQ 0(R12), DX                                        // DX = y[0]
    MULXQ R13, CX, BX                                       // t[0], t[1] = y[0] * x[0]
    MULXQ R14, AX, BP
    ADOXQ AX, BX
    MULXQ R15, AX, SI
    ADOXQ AX, BP
    MULXQ 24(R9), AX, DI
    ADOXQ AX, SI
    MULXQ 32(R9), AX, R8
    ADOXQ AX, DI
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 1
    XORQ DX, DX                                            // clear up flags
    MOVQ 8(R12), DX                                        // DX = y[1]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 2
    XORQ DX, DX                                            // clear up flags
    MOVQ 16(R12), DX                                       // DX = y[2]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 3
    XORQ DX, DX                                            // clear up flags
    MOVQ 24(R12), DX                                       // DX = y[3]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 4
    XORQ DX, DX                                            // clear up flags
    MOVQ 32(R12), DX                                       // DX = y[4]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    // outter loop 5
    XORQ DX, DX                                            // clear up flags
    MOVQ 40(R12), DX                                       // DX = y[5]
    MULXQ R13, AX, R11
    ADOXQ AX, CX
    ADCXQ R11, BX                                           // t[1] += regA
    MULXQ R14, AX, R11
    ADOXQ AX, BX
    ADCXQ R11, BP                                           // t[2] += regA
    MULXQ R15, AX, R11
    ADOXQ AX, BP
    ADCXQ R11, SI                                           // t[3] += regA
    MULXQ 24(R9), AX, R11
    ADOXQ AX, SI
    ADCXQ R11, DI                                           // t[4] += regA
    MULXQ 32(R9), AX, R11
    ADOXQ AX, DI
    ADCXQ R11, R8                                           // t[5] += regA
    MULXQ 40(R9), AX, R11
    ADOXQ AX, R8
    // add the last carries to R11
    MOVQ $0x0000000000000000, DX
    ADCXQ DX, R11
    ADOXQ DX, R11
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R11, R8
    MOVQ res+0(FP), R9                                     // dereference res
    MOVQ CX, R10
    SUBQ ·modulusElement+0(SB), R10
    MOVQ BX, R12
    SBBQ ·modulusElement+8(SB), R12
    MOVQ BP, R11
    SBBQ ·modulusElement+16(SB), R11
    MOVQ SI, R13
    SBBQ ·modulusElement+24(SB), R13
    MOVQ DI, R14
    SBBQ ·modulusElement+32(SB), R14
    MOVQ R8, R15
    SBBQ ·modulusElement+40(SB), R15
    CMOVQCC R10, CX
    CMOVQCC R12, BX
    CMOVQCC R11, BP
    CMOVQCC R13, SI
    CMOVQCC R14, DI
    CMOVQCC R15, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET
no_adx:
    MOVQ y+16(FP), R14                                     // dereference y
    MOVQ 0(R9), AX
    MOVQ 0(R14), R12
    MULQ R12
    MOVQ AX, CX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    MOVQ R11, BX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    MOVQ R11, BP
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    MOVQ R11, SI
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    MOVQ R11, DI
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    MOVQ R11, R8
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 8(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 16(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 24(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 32(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ 0(R9), AX
    MOVQ 40(R14), R12
    MULQ R12
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x8508bfffffffffff, R13
    IMULQ CX, R13
    MOVQ $0x8508c00000000001, AX
    MULQ R13
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R10
    MOVQ 8(R9), AX
    MULQ R12
    ADDQ R11, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R13
    ADDQ BX, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, CX
    MOVQ DX, R10
    MOVQ 16(R9), AX
    MULQ R12
    ADDQ R11, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R13
    ADDQ BP, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BX
    MOVQ DX, R10
    MOVQ 24(R9), AX
    MULQ R12
    ADDQ R11, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R13
    ADDQ SI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, BP
    MOVQ DX, R10
    MOVQ 32(R9), AX
    MULQ R12
    ADDQ R11, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R13
    ADDQ DI, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, SI
    MOVQ DX, R10
    MOVQ 40(R9), AX
    MULQ R12
    ADDQ R11, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R13
    ADDQ R8, R10
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R10
    ADCQ $0x0000000000000000, DX
    MOVQ R10, DI
    MOVQ DX, R10
    ADDQ R10, R11
    MOVQ R11, R8
    MOVQ res+0(FP), R9                                     // dereference res
    MOVQ CX, R15
    SUBQ ·modulusElement+0(SB), R15
    MOVQ BX, R10
    SBBQ ·modulusElement+8(SB), R10
    MOVQ BP, R12
    SBBQ ·modulusElement+16(SB), R12
    MOVQ SI, R11
    SBBQ ·modulusElement+24(SB), R11
    MOVQ DI, R13
    SBBQ ·modulusElement+32(SB), R13
    MOVQ R8, R14
    SBBQ ·modulusElement+40(SB), R14
    CMOVQCC R15, CX
    CMOVQCC R10, BX
    CMOVQCC R12, BP
    CMOVQCC R11, SI
    CMOVQCC R13, DI
    CMOVQCC R14, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
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
    MOVQ res+0(FP), R9                                     // dereference x
    MOVQ 0(R9), CX                                         // t[0] = x[0]
    MOVQ 8(R9), BX                                         // t[1] = x[1]
    MOVQ 16(R9), BP                                        // t[2] = x[2]
    MOVQ 24(R9), SI                                        // t[3] = x[3]
    MOVQ 32(R9), DI                                        // t[4] = x[4]
    MOVQ 40(R9), R8                                        // t[5] = x[5]
    CMPB ·supportAdx(SB), $0x0000000000000001             // check if we support MULX and ADOX instructions
    JNE no_adx                                            // no support for MULX or ADOX instructions
    // outter loop 0
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ AX, R8
    // outter loop 1
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ AX, R8
    // outter loop 2
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ AX, R8
    // outter loop 3
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ AX, R8
    // outter loop 4
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ AX, R8
    // outter loop 5
    XORQ DX, DX                                            // clear up flags
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear the flags
    // C,_ := t[0] + m*q[0]
    MULXQ ·modulusElement+0(SB), AX, R10
    ADCXQ CX, AX
    MOVQ R10, CX
    // for j=1 to N-1
    //     (C,t[j-1]) := t[j] + m*q[j] + C
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ AX, R8
    MOVQ CX, R11
    SUBQ ·modulusElement+0(SB), R11
    MOVQ BX, R12
    SBBQ ·modulusElement+8(SB), R12
    MOVQ BP, R13
    SBBQ ·modulusElement+16(SB), R13
    MOVQ SI, R14
    SBBQ ·modulusElement+24(SB), R14
    MOVQ DI, R15
    SBBQ ·modulusElement+32(SB), R15
    MOVQ R8, R10
    SBBQ ·modulusElement+40(SB), R10
    CMOVQCC R11, CX
    CMOVQCC R12, BX
    CMOVQCC R13, BP
    CMOVQCC R14, SI
    CMOVQCC R15, DI
    CMOVQCC R10, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET
no_adx:
    MOVQ $0x8508bfffffffffff, R14
    IMULQ CX, R14
    MOVQ $0x8508c00000000001, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R14
    ADDQ DI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, SI
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R14
    ADDQ R8, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, DI
    MOVQ DX, R11
    MOVQ R11, R8
    MOVQ $0x8508bfffffffffff, R14
    IMULQ CX, R14
    MOVQ $0x8508c00000000001, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R14
    ADDQ DI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, SI
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R14
    ADDQ R8, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, DI
    MOVQ DX, R11
    MOVQ R11, R8
    MOVQ $0x8508bfffffffffff, R14
    IMULQ CX, R14
    MOVQ $0x8508c00000000001, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R14
    ADDQ DI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, SI
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R14
    ADDQ R8, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, DI
    MOVQ DX, R11
    MOVQ R11, R8
    MOVQ $0x8508bfffffffffff, R14
    IMULQ CX, R14
    MOVQ $0x8508c00000000001, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R14
    ADDQ DI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, SI
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R14
    ADDQ R8, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, DI
    MOVQ DX, R11
    MOVQ R11, R8
    MOVQ $0x8508bfffffffffff, R14
    IMULQ CX, R14
    MOVQ $0x8508c00000000001, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R14
    ADDQ DI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, SI
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R14
    ADDQ R8, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, DI
    MOVQ DX, R11
    MOVQ R11, R8
    MOVQ $0x8508bfffffffffff, R14
    IMULQ CX, R14
    MOVQ $0x8508c00000000001, AX
    MULQ R14
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R11
    MOVQ $0x170b5d4430000000, AX
    MULQ R14
    ADDQ BX, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, CX
    MOVQ DX, R11
    MOVQ $0x1ef3622fba094800, AX
    MULQ R14
    ADDQ BP, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BX
    MOVQ DX, R11
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R14
    ADDQ SI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, BP
    MOVQ DX, R11
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R14
    ADDQ DI, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, SI
    MOVQ DX, R11
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R14
    ADDQ R8, R11
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R11
    ADCQ $0x0000000000000000, DX
    MOVQ R11, DI
    MOVQ DX, R11
    MOVQ R11, R8
    MOVQ CX, R10
    SUBQ ·modulusElement+0(SB), R10
    MOVQ BX, R11
    SBBQ ·modulusElement+8(SB), R11
    MOVQ BP, R12
    SBBQ ·modulusElement+16(SB), R12
    MOVQ SI, R13
    SBBQ ·modulusElement+24(SB), R13
    MOVQ DI, R14
    SBBQ ·modulusElement+32(SB), R14
    MOVQ R8, R15
    SBBQ ·modulusElement+40(SB), R15
    CMOVQCC R10, CX
    CMOVQCC R11, BX
    CMOVQCC R12, BP
    CMOVQCC R13, SI
    CMOVQCC R14, DI
    CMOVQCC R15, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
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
    MOVQ y+8(FP), R9                                       // dereference y
    // outter loop 0
    XORQ AX, AX                                            // clear up flags
    // dx = y[0]
    MOVQ 0(R9), DX
    MULXQ 8(R9), R11, R12
    MULXQ 16(R9), AX, R13
    ADCXQ AX, R12
    MULXQ 24(R9), AX, R14
    ADCXQ AX, R13
    MULXQ 32(R9), AX, R15
    ADCXQ AX, R14
    MULXQ 40(R9), AX, R10
    ADCXQ AX, R15
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, CX, DX
    ADCXQ R11, R11
    MOVQ R11, BX
    ADOXQ DX, BX
    ADCXQ R12, R12
    MOVQ R12, BP
    ADOXQ AX, BP
    ADCXQ R13, R13
    MOVQ R13, SI
    ADOXQ AX, SI
    ADCXQ R14, R14
    MOVQ R14, DI
    ADOXQ AX, DI
    ADCXQ R15, R15
    MOVQ R15, R8
    ADOXQ AX, R8
    ADCXQ R10, R10
    ADOXQ AX, R10
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·modulusElement+0(SB), AX, R11
    ADCXQ CX, AX
    MOVQ R11, CX
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R10, R8
    // outter loop 1
    XORQ AX, AX                                            // clear up flags
    // dx = y[1]
    MOVQ 8(R9), DX
    MULXQ 16(R9), R12, R13
    MULXQ 24(R9), AX, R14
    ADCXQ AX, R13
    MULXQ 32(R9), AX, R15
    ADCXQ AX, R14
    MULXQ 40(R9), AX, R10
    ADCXQ AX, R15
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    ADCXQ R12, R12
    ADOXQ R12, BP
    ADCXQ R13, R13
    ADOXQ R13, SI
    ADCXQ R14, R14
    ADOXQ R14, DI
    ADCXQ R15, R15
    ADOXQ R15, R8
    ADCXQ R10, R10
    ADOXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, AX, DX
    ADOXQ AX, BX
    MOVQ $0x0000000000000000, AX
    ADOXQ DX, BP
    ADOXQ AX, SI
    ADOXQ AX, DI
    ADOXQ AX, R8
    ADOXQ AX, R10
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·modulusElement+0(SB), AX, R11
    ADCXQ CX, AX
    MOVQ R11, CX
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R10, R8
    // outter loop 2
    XORQ AX, AX                                            // clear up flags
    // dx = y[2]
    MOVQ 16(R9), DX
    MULXQ 24(R9), R12, R13
    MULXQ 32(R9), AX, R14
    ADCXQ AX, R13
    MULXQ 40(R9), AX, R10
    ADCXQ AX, R14
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    ADCXQ R12, R12
    ADOXQ R12, SI
    ADCXQ R13, R13
    ADOXQ R13, DI
    ADCXQ R14, R14
    ADOXQ R14, R8
    ADCXQ R10, R10
    ADOXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, AX, DX
    ADOXQ AX, BP
    MOVQ $0x0000000000000000, AX
    ADOXQ DX, SI
    ADOXQ AX, DI
    ADOXQ AX, R8
    ADOXQ AX, R10
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·modulusElement+0(SB), AX, R15
    ADCXQ CX, AX
    MOVQ R15, CX
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R10, R8
    // outter loop 3
    XORQ AX, AX                                            // clear up flags
    // dx = y[3]
    MOVQ 24(R9), DX
    MULXQ 32(R9), R11, R12
    MULXQ 40(R9), AX, R10
    ADCXQ AX, R12
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    ADCXQ R11, R11
    ADOXQ R11, DI
    ADCXQ R12, R12
    ADOXQ R12, R8
    ADCXQ R10, R10
    ADOXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, AX, DX
    ADOXQ AX, SI
    MOVQ $0x0000000000000000, AX
    ADOXQ DX, DI
    ADOXQ AX, R8
    ADOXQ AX, R10
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·modulusElement+0(SB), AX, R13
    ADCXQ CX, AX
    MOVQ R13, CX
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R10, R8
    // outter loop 4
    XORQ AX, AX                                            // clear up flags
    // dx = y[4]
    MOVQ 32(R9), DX
    MULXQ 40(R9), R14, R10
    ADCXQ R14, R14
    ADOXQ R14, R8
    ADCXQ R10, R10
    ADOXQ AX, R10
    XORQ AX, AX                                            // clear up flags
    MULXQ DX, AX, DX
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADOXQ DX, R8
    ADOXQ AX, R10
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·modulusElement+0(SB), AX, R15
    ADCXQ CX, AX
    MOVQ R15, CX
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R10, R8
    // outter loop 5
    XORQ AX, AX                                            // clear up flags
    // dx = y[5]
    MOVQ 40(R9), DX
    MULXQ DX, AX, R10
    ADCXQ AX, R8
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R10
    MOVQ CX, DX
    MULXQ ·modulusElementInv0(SB), DX, AX                  // m := t[0]*q'[0] mod W
    XORQ AX, AX                                            // clear up flags
    MULXQ ·modulusElement+0(SB), AX, R11
    ADCXQ CX, AX
    MOVQ R11, CX
    ADCXQ BX, CX
    MULXQ ·modulusElement+8(SB), AX, BX
    ADOXQ AX, CX
    ADCXQ BP, BX
    MULXQ ·modulusElement+16(SB), AX, BP
    ADOXQ AX, BX
    ADCXQ SI, BP
    MULXQ ·modulusElement+24(SB), AX, SI
    ADOXQ AX, BP
    ADCXQ DI, SI
    MULXQ ·modulusElement+32(SB), AX, DI
    ADOXQ AX, SI
    ADCXQ R8, DI
    MULXQ ·modulusElement+40(SB), AX, R8
    ADOXQ AX, DI
    MOVQ $0x0000000000000000, AX
    ADCXQ AX, R8
    ADOXQ R10, R8
    // dereference res
    MOVQ res+0(FP), R12
    MOVQ CX, R13
    SUBQ ·modulusElement+0(SB), R13
    MOVQ BX, R14
    SBBQ ·modulusElement+8(SB), R14
    MOVQ BP, R15
    SBBQ ·modulusElement+16(SB), R15
    MOVQ SI, R11
    SBBQ ·modulusElement+24(SB), R11
    MOVQ DI, R9
    SBBQ ·modulusElement+32(SB), R9
    MOVQ R8, R10
    SBBQ ·modulusElement+40(SB), R10
    CMOVQCC R13, CX
    CMOVQCC R14, BX
    CMOVQCC R15, BP
    CMOVQCC R11, SI
    CMOVQCC R9, DI
    CMOVQCC R10, R8
    MOVQ CX, 0(R12)
    MOVQ BX, 8(R12)
    MOVQ BP, 16(R12)
    MOVQ SI, 24(R12)
    MOVQ DI, 32(R12)
    MOVQ R8, 40(R12)
    RET
no_adx:
    // dereference y
    MOVQ y+8(FP), R9
    MOVQ 0(R9), AX
    MOVQ 0(R9), R14
    MULQ R14
    MOVQ AX, CX
    MOVQ DX, R15
    MOVQ $0x8508bfffffffffff, R11
    IMULQ CX, R11
    MOVQ $0x8508c00000000001, AX
    MULQ R11
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ 8(R9), AX
    MULQ R14
    MOVQ R15, BX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x170b5d4430000000, AX
    MULQ R11
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ 16(R9), AX
    MULQ R14
    MOVQ R15, BP
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1ef3622fba094800, AX
    MULQ R11
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ 24(R9), AX
    MULQ R14
    MOVQ R15, SI
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R11
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ 32(R9), AX
    MULQ R14
    MOVQ R15, DI
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R11
    ADDQ DI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, SI
    MOVQ DX, R13
    MOVQ 40(R9), AX
    MULQ R14
    MOVQ R15, R8
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R11
    ADDQ R8, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, DI
    MOVQ DX, R13
    ADDQ R13, R15
    MOVQ R15, R8
    MOVQ 0(R9), AX
    MOVQ 8(R9), R14
    MULQ R14
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x8508bfffffffffff, R11
    IMULQ CX, R11
    MOVQ $0x8508c00000000001, AX
    MULQ R11
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ 8(R9), AX
    MULQ R14
    ADDQ R15, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x170b5d4430000000, AX
    MULQ R11
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ 16(R9), AX
    MULQ R14
    ADDQ R15, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1ef3622fba094800, AX
    MULQ R11
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ 24(R9), AX
    MULQ R14
    ADDQ R15, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R11
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ 32(R9), AX
    MULQ R14
    ADDQ R15, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R11
    ADDQ DI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, SI
    MOVQ DX, R13
    MOVQ 40(R9), AX
    MULQ R14
    ADDQ R15, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R11
    ADDQ R8, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, DI
    MOVQ DX, R13
    ADDQ R13, R15
    MOVQ R15, R8
    MOVQ 0(R9), AX
    MOVQ 16(R9), R14
    MULQ R14
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x8508bfffffffffff, R11
    IMULQ CX, R11
    MOVQ $0x8508c00000000001, AX
    MULQ R11
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ 8(R9), AX
    MULQ R14
    ADDQ R15, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x170b5d4430000000, AX
    MULQ R11
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ 16(R9), AX
    MULQ R14
    ADDQ R15, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1ef3622fba094800, AX
    MULQ R11
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ 24(R9), AX
    MULQ R14
    ADDQ R15, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R11
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ 32(R9), AX
    MULQ R14
    ADDQ R15, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R11
    ADDQ DI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, SI
    MOVQ DX, R13
    MOVQ 40(R9), AX
    MULQ R14
    ADDQ R15, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R11
    ADDQ R8, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, DI
    MOVQ DX, R13
    ADDQ R13, R15
    MOVQ R15, R8
    MOVQ 0(R9), AX
    MOVQ 24(R9), R14
    MULQ R14
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x8508bfffffffffff, R11
    IMULQ CX, R11
    MOVQ $0x8508c00000000001, AX
    MULQ R11
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ 8(R9), AX
    MULQ R14
    ADDQ R15, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x170b5d4430000000, AX
    MULQ R11
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ 16(R9), AX
    MULQ R14
    ADDQ R15, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1ef3622fba094800, AX
    MULQ R11
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ 24(R9), AX
    MULQ R14
    ADDQ R15, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R11
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ 32(R9), AX
    MULQ R14
    ADDQ R15, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R11
    ADDQ DI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, SI
    MOVQ DX, R13
    MOVQ 40(R9), AX
    MULQ R14
    ADDQ R15, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R11
    ADDQ R8, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, DI
    MOVQ DX, R13
    ADDQ R13, R15
    MOVQ R15, R8
    MOVQ 0(R9), AX
    MOVQ 32(R9), R14
    MULQ R14
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x8508bfffffffffff, R11
    IMULQ CX, R11
    MOVQ $0x8508c00000000001, AX
    MULQ R11
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ 8(R9), AX
    MULQ R14
    ADDQ R15, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x170b5d4430000000, AX
    MULQ R11
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ 16(R9), AX
    MULQ R14
    ADDQ R15, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1ef3622fba094800, AX
    MULQ R11
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ 24(R9), AX
    MULQ R14
    ADDQ R15, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R11
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ 32(R9), AX
    MULQ R14
    ADDQ R15, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R11
    ADDQ DI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, SI
    MOVQ DX, R13
    MOVQ 40(R9), AX
    MULQ R14
    ADDQ R15, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R11
    ADDQ R8, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, DI
    MOVQ DX, R13
    ADDQ R13, R15
    MOVQ R15, R8
    MOVQ 0(R9), AX
    MOVQ 40(R9), R14
    MULQ R14
    ADDQ AX, CX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x8508bfffffffffff, R11
    IMULQ CX, R11
    MOVQ $0x8508c00000000001, AX
    MULQ R11
    ADDQ CX, AX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R13
    MOVQ 8(R9), AX
    MULQ R14
    ADDQ R15, BX
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BX
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x170b5d4430000000, AX
    MULQ R11
    ADDQ BX, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, CX
    MOVQ DX, R13
    MOVQ 16(R9), AX
    MULQ R14
    ADDQ R15, BP
    ADCQ $0x0000000000000000, DX
    ADDQ AX, BP
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1ef3622fba094800, AX
    MULQ R11
    ADDQ BP, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BX
    MOVQ DX, R13
    MOVQ 24(R9), AX
    MULQ R14
    ADDQ R15, SI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, SI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x1a22d9f300f5138f, AX
    MULQ R11
    ADDQ SI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, BP
    MOVQ DX, R13
    MOVQ 32(R9), AX
    MULQ R14
    ADDQ R15, DI
    ADCQ $0x0000000000000000, DX
    ADDQ AX, DI
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0xc63b05c06ca1493b, AX
    MULQ R11
    ADDQ DI, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, SI
    MOVQ DX, R13
    MOVQ 40(R9), AX
    MULQ R14
    ADDQ R15, R8
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R8
    ADCQ $0x0000000000000000, DX
    MOVQ DX, R15
    MOVQ $0x01ae3a4617c510ea, AX
    MULQ R11
    ADDQ R8, R13
    ADCQ $0x0000000000000000, DX
    ADDQ AX, R13
    ADCQ $0x0000000000000000, DX
    MOVQ R13, DI
    MOVQ DX, R13
    ADDQ R13, R15
    MOVQ R15, R8
    // dereference res
    MOVQ res+0(FP), R12
    MOVQ CX, R10
    SUBQ ·modulusElement+0(SB), R10
    MOVQ BX, R13
    SBBQ ·modulusElement+8(SB), R13
    MOVQ BP, R14
    SBBQ ·modulusElement+16(SB), R14
    MOVQ SI, R15
    SBBQ ·modulusElement+24(SB), R15
    MOVQ DI, R11
    SBBQ ·modulusElement+32(SB), R11
    MOVQ R8, R9
    SBBQ ·modulusElement+40(SB), R9
    CMOVQCC R10, CX
    CMOVQCC R13, BX
    CMOVQCC R14, BP
    CMOVQCC R15, SI
    CMOVQCC R11, DI
    CMOVQCC R9, R8
    MOVQ CX, 0(R12)
    MOVQ BX, 8(R12)
    MOVQ BP, 16(R12)
    MOVQ SI, 24(R12)
    MOVQ DI, 32(R12)
    MOVQ R8, 40(R12)
    RET

TEXT ·reduceElement(SB), NOSPLIT, $0-8
    MOVQ res+0(FP), R9                                     // dereference x
    MOVQ 0(R9), CX                                         // t[0] = x[0]
    MOVQ 8(R9), BX                                         // t[1] = x[1]
    MOVQ 16(R9), BP                                        // t[2] = x[2]
    MOVQ 24(R9), SI                                        // t[3] = x[3]
    MOVQ 32(R9), DI                                        // t[4] = x[4]
    MOVQ 40(R9), R8                                        // t[5] = x[5]
    MOVQ CX, R10
    SUBQ ·modulusElement+0(SB), R10
    MOVQ BX, R11
    SBBQ ·modulusElement+8(SB), R11
    MOVQ BP, R12
    SBBQ ·modulusElement+16(SB), R12
    MOVQ SI, R13
    SBBQ ·modulusElement+24(SB), R13
    MOVQ DI, R14
    SBBQ ·modulusElement+32(SB), R14
    MOVQ R8, R15
    SBBQ ·modulusElement+40(SB), R15
    CMOVQCC R10, CX
    CMOVQCC R11, BX
    CMOVQCC R12, BP
    CMOVQCC R13, SI
    CMOVQCC R14, DI
    CMOVQCC R15, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET

TEXT ·addElement(SB), NOSPLIT, $0-24
    MOVQ x+8(FP), R9                                       // dereference x
    MOVQ y+16(FP), R10                                     // dereference y
    MOVQ 0(R9), CX                                         // t[0] = x[0]
    MOVQ 8(R9), BX                                         // t[1] = x[1]
    MOVQ 16(R9), BP                                        // t[2] = x[2]
    MOVQ 24(R9), SI                                        // t[3] = x[3]
    MOVQ 32(R9), DI                                        // t[4] = x[4]
    MOVQ 40(R9), R8                                        // t[5] = x[5]
    ADDQ 0(R10), CX
    ADCQ 8(R10), BX
    ADCQ 16(R10), BP
    ADCQ 24(R10), SI
    ADCQ 32(R10), DI
    ADCQ 40(R10), R8
    MOVQ res+0(FP), R9                                     // dereference res
    MOVQ CX, R11
    SUBQ ·modulusElement+0(SB), R11
    MOVQ BX, R12
    SBBQ ·modulusElement+8(SB), R12
    MOVQ BP, R13
    SBBQ ·modulusElement+16(SB), R13
    MOVQ SI, R14
    SBBQ ·modulusElement+24(SB), R14
    MOVQ DI, R15
    SBBQ ·modulusElement+32(SB), R15
    MOVQ R8, R10
    SBBQ ·modulusElement+40(SB), R10
    CMOVQCC R11, CX
    CMOVQCC R12, BX
    CMOVQCC R13, BP
    CMOVQCC R14, SI
    CMOVQCC R15, DI
    CMOVQCC R10, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET

TEXT ·addAssignElement(SB), NOSPLIT, $0-16
    MOVQ res+0(FP), R9                                     // dereference x
    MOVQ y+8(FP), R10                                      // dereference y
    MOVQ 0(R9), CX                                         // t[0] = x[0]
    MOVQ 8(R9), BX                                         // t[1] = x[1]
    MOVQ 16(R9), BP                                        // t[2] = x[2]
    MOVQ 24(R9), SI                                        // t[3] = x[3]
    MOVQ 32(R9), DI                                        // t[4] = x[4]
    MOVQ 40(R9), R8                                        // t[5] = x[5]
    ADDQ 0(R10), CX
    ADCQ 8(R10), BX
    ADCQ 16(R10), BP
    ADCQ 24(R10), SI
    ADCQ 32(R10), DI
    ADCQ 40(R10), R8
    MOVQ CX, R11
    SUBQ ·modulusElement+0(SB), R11
    MOVQ BX, R12
    SBBQ ·modulusElement+8(SB), R12
    MOVQ BP, R13
    SBBQ ·modulusElement+16(SB), R13
    MOVQ SI, R14
    SBBQ ·modulusElement+24(SB), R14
    MOVQ DI, R15
    SBBQ ·modulusElement+32(SB), R15
    MOVQ R8, R10
    SBBQ ·modulusElement+40(SB), R10
    CMOVQCC R11, CX
    CMOVQCC R12, BX
    CMOVQCC R13, BP
    CMOVQCC R14, SI
    CMOVQCC R15, DI
    CMOVQCC R10, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET

TEXT ·doubleElement(SB), NOSPLIT, $0-16
    MOVQ res+0(FP), R9                                     // dereference x
    MOVQ y+8(FP), R10                                      // dereference y
    MOVQ 0(R10), CX                                        // t[0] = y[0]
    MOVQ 8(R10), BX                                        // t[1] = y[1]
    MOVQ 16(R10), BP                                       // t[2] = y[2]
    MOVQ 24(R10), SI                                       // t[3] = y[3]
    MOVQ 32(R10), DI                                       // t[4] = y[4]
    MOVQ 40(R10), R8                                       // t[5] = y[5]
    ADDQ CX, CX
    ADCQ BX, BX
    ADCQ BP, BP
    ADCQ SI, SI
    ADCQ DI, DI
    ADCQ R8, R8
    MOVQ CX, R11
    SUBQ ·modulusElement+0(SB), R11
    MOVQ BX, R12
    SBBQ ·modulusElement+8(SB), R12
    MOVQ BP, R13
    SBBQ ·modulusElement+16(SB), R13
    MOVQ SI, R14
    SBBQ ·modulusElement+24(SB), R14
    MOVQ DI, R15
    SBBQ ·modulusElement+32(SB), R15
    MOVQ R8, R10
    SBBQ ·modulusElement+40(SB), R10
    CMOVQCC R11, CX
    CMOVQCC R12, BX
    CMOVQCC R13, BP
    CMOVQCC R14, SI
    CMOVQCC R15, DI
    CMOVQCC R10, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET

TEXT ·subElement(SB), NOSPLIT, $0-24
    MOVQ x+8(FP), R9                                       // dereference x
    MOVQ y+16(FP), R10                                     // dereference y
    MOVQ 0(R9), CX                                         // t[0] = x[0]
    MOVQ 8(R9), BX                                         // t[1] = x[1]
    MOVQ 16(R9), BP                                        // t[2] = x[2]
    MOVQ 24(R9), SI                                        // t[3] = x[3]
    MOVQ 32(R9), DI                                        // t[4] = x[4]
    MOVQ 40(R9), R8                                        // t[5] = x[5]
    XORQ DX, DX
    SUBQ 0(R10), CX
    SBBQ 8(R10), BX
    SBBQ 16(R10), BP
    SBBQ 24(R10), SI
    SBBQ 32(R10), DI
    SBBQ 40(R10), R8
    MOVQ $0x8508c00000000001, R11
    MOVQ $0x170b5d4430000000, R12
    MOVQ $0x1ef3622fba094800, R13
    MOVQ $0x1a22d9f300f5138f, R14
    MOVQ $0xc63b05c06ca1493b, R15
    MOVQ $0x01ae3a4617c510ea, R10
    CMOVQCC DX, R11
    CMOVQCC DX, R12
    CMOVQCC DX, R13
    CMOVQCC DX, R14
    CMOVQCC DX, R15
    CMOVQCC DX, R10
    ADDQ R11, CX
    ADCQ R12, BX
    ADCQ R13, BP
    ADCQ R14, SI
    ADCQ R15, DI
    ADCQ R10, R8
    MOVQ res+0(FP), R9                                     // dereference res
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET

TEXT ·subAssignElement(SB), NOSPLIT, $0-16
    MOVQ res+0(FP), R9                                     // dereference x
    MOVQ y+8(FP), R10                                      // dereference y
    MOVQ 0(R9), CX                                         // t[0] = x[0]
    MOVQ 8(R9), BX                                         // t[1] = x[1]
    MOVQ 16(R9), BP                                        // t[2] = x[2]
    MOVQ 24(R9), SI                                        // t[3] = x[3]
    MOVQ 32(R9), DI                                        // t[4] = x[4]
    MOVQ 40(R9), R8                                        // t[5] = x[5]
    XORQ DX, DX
    SUBQ 0(R10), CX
    SBBQ 8(R10), BX
    SBBQ 16(R10), BP
    SBBQ 24(R10), SI
    SBBQ 32(R10), DI
    SBBQ 40(R10), R8
    MOVQ $0x8508c00000000001, R11
    MOVQ $0x170b5d4430000000, R12
    MOVQ $0x1ef3622fba094800, R13
    MOVQ $0x1a22d9f300f5138f, R14
    MOVQ $0xc63b05c06ca1493b, R15
    MOVQ $0x01ae3a4617c510ea, R10
    CMOVQCC DX, R11
    CMOVQCC DX, R12
    CMOVQCC DX, R13
    CMOVQCC DX, R14
    CMOVQCC DX, R15
    CMOVQCC DX, R10
    ADDQ R11, CX
    ADCQ R12, BX
    ADCQ R13, BP
    ADCQ R14, SI
    ADCQ R15, DI
    ADCQ R10, R8
    MOVQ CX, 0(R9)
    MOVQ BX, 8(R9)
    MOVQ BP, 16(R9)
    MOVQ SI, 24(R9)
    MOVQ DI, 32(R9)
    MOVQ R8, 40(R9)
    RET
