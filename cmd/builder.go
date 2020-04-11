package cmd

import (
	"fmt"
	"io"
)

const dx = "DX"
const ax = "AX"

type register string

func (r *register) at(wordOffset int) string {
	return fmt.Sprintf("%d(%s)", wordOffset*8, string(*r))
}

type asmBuilder struct {
	writer    io.Writer
	registers []register
}

func newAsmBuilder(w io.Writer) *asmBuilder {
	b := &asmBuilder{
		writer:    w,
		registers: make([]register, len(staticRegisters)),
	}
	copy(b.registers, staticRegisters)
	return b
}

func (builder *asmBuilder) PopRegister() register {
	r := builder.registers[0]
	builder.registers = builder.registers[1:]
	return r
}

func (builder *asmBuilder) PushRegister(r ...register) {
	builder.registers = append(builder.registers, r...)
}

func (builder *asmBuilder) Comment(s string) {
	builder.WriteLn("    // " + s)
}

func (builder *asmBuilder) WriteLn(s string) {
	builder.Write(s + "\n")
}

func (builder *asmBuilder) Write(s string) {
	builder.writer.Write([]byte(s))
}

func (builder *asmBuilder) RET() {
	builder.WriteLn("    RET")
}

func (builder *asmBuilder) MULXQ(src, lo, hi interface{}) {
	builder.writeOp("MULXQ", src, lo, hi)
}

func (builder *asmBuilder) SUBQ(r1, r2 interface{}) {
	builder.writeOp("SUBQ", r1, r2)
}

func (builder *asmBuilder) SBBQ(r1, r2 interface{}) {
	builder.writeOp("SBBQ", r1, r2)
}

func (builder *asmBuilder) ADDQ(r1, r2 interface{}) {
	builder.writeOp("ADDQ", r1, r2)
}

func (builder *asmBuilder) ADCQ(r1, r2 interface{}) {
	builder.writeOp("ADCQ", r1, r2)
}

func (builder *asmBuilder) ADOXQ(r1, r2 interface{}) {
	builder.writeOp("ADOXQ", r1, r2)
}

func (builder *asmBuilder) ADCXQ(r1, r2 interface{}) {
	builder.writeOp("ADCXQ", r1, r2)
}

func (builder *asmBuilder) XORQ(r1, r2 interface{}) {
	builder.writeOp("XORQ", r1, r2)
}

func (builder *asmBuilder) MOVQ(r1, r2 interface{}) {
	builder.writeOp("MOVQ", r1, r2)
}

func (builder *asmBuilder) IMULQ(r1, r2 interface{}) {
	builder.writeOp("IMULQ", r1, r2)
}

func (builder *asmBuilder) MULQ(r1 interface{}) {
	builder.writeOp("MULQ", r1)
}

func (builder *asmBuilder) CMPB(r1, r2 interface{}) {
	builder.writeOp("CMPB", r1, r2)
}

func (builder *asmBuilder) JNE(label string) {
	builder.writeOp("JNE", label)
}

func (builder *asmBuilder) JCS(label string) {
	builder.writeOp("JCS", label)
}

func (builder *asmBuilder) JMP(label string) {
	builder.writeOp("JMP", label)
}

func (builder *asmBuilder) writeOp(instruction string, r0 interface{}, r ...interface{}) {
	builder.Write(fmt.Sprintf("    %s %s", instruction, op(r0)))
	for _, rn := range r {
		builder.Write(fmt.Sprintf(", %s", op(rn)))
	}
	builder.Write("\n")
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

var staticRegisters = []register{ // AX and DX are reserved
	"CX",
	"BX",
	"BP",
	"SI",
	"DI",
	"R8",
	"R9",
	"R10",
	"R11",
	"R12",
	"R13",
	"R14",
	"R15",
}
