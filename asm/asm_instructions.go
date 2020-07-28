package asm

import "fmt"

func ret() {
	writeLn("    RET")
}

func mulxq(src, lo, hi interface{}, comment ...string) {
	writeOp(comment, "MULXQ", src, lo, hi)
}

func subq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "SUBQ", r1, r2)
}

func sbbq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "SBBQ", r1, r2)
}

func addq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "ADDQ", r1, r2)
}

func adcq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "ADCQ", r1, r2)
}

func adoxq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "ADOXQ", r1, r2)
}

func adcxq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "ADCXQ", r1, r2)
}

func xorq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "XORQ", r1, r2)
}

func xorps(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "XORPS", r1, r2)
}

func movq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "MOVQ", r1, r2)
}

func movups(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "MOVUPS", r1, r2)
}

func movntiq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "MOVNTIQ", r1, r2)
}

func pushq(r1 interface{}, comment ...string) {
	writeOp(comment, "PUSHQ", r1)
}

func popq(r1 interface{}, comment ...string) {
	writeOp(comment, "POPQ", r1)
}

func imulq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "IMULQ", r1, r2)
}

func mulq(r1 interface{}, comment ...string) {
	writeOp(comment, "MULQ", r1)
}

func cmpb(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "CMPB", r1, r2)
}

func cmpq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "CMPQ", r1, r2)
}

func orq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "ORQ", r1, r2)
}

func testq(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "TESTQ", r1, r2)
}

func cmovqcc(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "CMOVQCC", r1, r2)
}

func cmovqcs(r1, r2 interface{}, comment ...string) {
	writeOp(comment, "CMOVQCS", r1, r2)
}

func label(l lbl) {
	writeLn(string(l) + ":")
}

func jne(label lbl, comment ...string) {
	writeOp(comment, "JNE", string(label))
}

func jcs(label lbl, comment ...string) {
	writeOp(comment, "JCS", string(label))
}

func jcc(label lbl, comment ...string) {
	writeOp(comment, "JCC", string(label))
}

func jmp(label lbl, comment ...string) {
	writeOp(comment, "JMP", string(label))
}

func (builder *assembly) reset() {
	builder.registers = make([]register, len(staticRegisters))
	copy(builder.registers, staticRegisters)
}

func comment(s string) {
	writeLn("    // " + s)
}

func fnHeader(funcName string, stackSize, argSize int, reserved ...register) {
	writeLn("")
	var header string
	if stackSize == 0 {
		header = "TEXT ·%s(SB), NOSPLIT, $%d-%d"
	} else {
		header = "TEXT ·%s(SB), $%d-%d"
	}

	writeLn(fmt.Sprintf(header, funcName, stackSize, argSize))
	builder.reset()
	for _, r := range reserved {
		builder.remove(r)
	}

}

func (b *assembly) remove(r register) {
	for j := 0; j < len(builder.registers); j++ {
		if builder.registers[j] == r {
			builder.registers[j] = builder.registers[len(builder.registers)-1]
			builder.registers = builder.registers[:len(builder.registers)-1]
			return
		}
	}
	panic("register not found")
}

func writeLn(s string) {
	write(s + "\n")
}

func write(s string) {
	builder.writer.Write([]byte(s))
}

func writeOp(comments []string, instruction string, r0 interface{}, r ...interface{}) {
	write(fmt.Sprintf("    %s %s", instruction, op(r0)))
	l := len(op(r0))
	for _, rn := range r {
		write(fmt.Sprintf(", %s", op(rn)))
		l += (2 + len(op(rn)))
	}
	if len(comments) == 1 {
		l = 50 - l
		for i := 0; i < l; i++ {
			write(" ")
		}
		write("// " + comments[0])
	}
	write("\n")
}

func op(i interface{}) string {
	switch t := i.(type) {
	case string:
		return t
	case register:
		return string(t)
	case int:
		return fmt.Sprintf("$%#016x", uint64(t))
	case uint64:
		return fmt.Sprintf("$%#016x", t)
	}
	panic("unsupported interface type")
}
