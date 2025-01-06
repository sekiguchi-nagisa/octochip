package main

import (
	"fmt"
	"io"
)

var instructionBuilders = map[InstructionType]func(uint16, InstructionType, byte, byte, byte) Instruction{
	OP_0NNN: NewAddrIns,
	OP_00E0: NewZeroIns,
	OP_00EE: NewZeroIns,
	OP_1NNN: NewAddrIns,
	OP_2NNN: NewAddrIns,
	OP_3XNN: NewOneRegConstIns,
	OP_4XNN: NewOneRegConstIns,
	OP_5XY0: NewTwoRegIns,
	OP_6XNN: NewOneRegConstIns,
	OP_7XNN: NewOneRegConstIns,
	OP_8XY0: NewTwoRegIns,
	OP_8XY1: NewTwoRegIns,
	OP_8XY2: NewTwoRegIns,
	OP_8XY3: NewTwoRegIns,
	OP_8XY4: NewTwoRegIns,
	OP_8XY5: NewTwoRegIns,
	OP_8XY6: NewTwoRegIns,
	OP_8XY7: NewTwoRegIns,
	OP_8XYE: NewTwoRegIns,
	OP_9XY0: NewTwoRegIns,
	OP_ANNN: NewAddrIns,
	OP_BNNN: NewAddrIns,
	OP_CXNN: NewOneRegConstIns,
	OP_DXYN: NewTwoRegConstIns,
	OP_EX9E: NewOneRegIns,
	OP_EXA1: NewOneRegIns,
	OP_FX07: NewOneRegIns,
	OP_FX0A: NewOneRegIns,
	OP_FX15: NewOneRegIns,
	OP_FX18: NewOneRegIns,
	OP_FX1E: NewOneRegIns,
	OP_FX29: NewOneRegIns,
	OP_FX33: NewOneRegIns,
	OP_FX55: NewOneRegIns,
	OP_FX65: NewOneRegIns,
}

func Disassemble(reader io.Reader, writer io.Writer) error {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// generate instruction sequence
	var instructionSeq []Instruction
	for i := 0; i < len(buf); i++ {
		if i+1 < len(buf) {
			op, r1, r2, r3 := DecodeInstruction(buf[i], buf[i+1])
			addr := uint16(Chip8ProgStartAddr + i)
			var ins Instruction
			if op != OP_INVALID {
				ins = instructionBuilders[op](addr, op, r1, r2, r3)
			} else {
				ins = InvalidIns{addr: addr, b1: buf[i], b2: buf[i+1]}
			}
			instructionSeq = append(instructionSeq, ins)
			i++
		}
	}

	// resolve jump target
	labelMap := make(map[uint16]string)
	labelIdCount := 0
	for _, ins := range instructionSeq {
		switch ins := ins.(type) {
		case AddrIns:
			prefix := ""
			if ins.op == OP_1NNN {
				prefix = "label"
			} else if ins.op == OP_2NNN {
				prefix = "subroutine"
			} else {
				continue
			}
			labelMap[ins.target] = fmt.Sprintf("%s%d", prefix, labelIdCount)
			labelIdCount++
		}
	}

	// print sequence
	printer := InstructionPrinter{labelMap: labelMap, writer: writer}
	_, _ = fmt.Fprintln(printer.writer, "start:")
	for _, ins := range instructionSeq {
		if label, ok := labelMap[ins.Address()]; ok {
			_, _ = fmt.Fprintf(printer.writer, "%s:\n", label)
		}
		if err := ins.Print(printer); err != nil {
			return err
		}
	}
	return nil
}
