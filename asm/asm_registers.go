package asm

import "fmt"

const dx = register("DX")
const ax = register("AX")

type lbl string
type register string

func (r *register) at(wordOffset int) string {
	return fmt.Sprintf("%d(%s)", wordOffset*8, string(*r))
}

func availableRegisters() int {
	return len(builder.registers)
}

func popRegister() register {
	r := builder.registers[0]
	builder.registers = builder.registers[1:]
	return r
}

func popRegisters(n int) []register {
	toReturn := make([]register, n)
	for i := 0; i < n; i++ {
		toReturn[i] = popRegister()
	}
	return toReturn
}

func pushRegister(r ...register) {
	builder.registers = append(builder.registers, r...)
}

var staticRegisters = []register{
	"AX",
	"DX",
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

var labelCounter = 0

func newLabel() lbl {
	labelCounter++
	return lbl(fmt.Sprintf("l%d", labelCounter))
}
