package main

import (
	"fmt"
	"io"
)

type InstructionPrinter struct {
	writer   io.Writer
	labelMap map[uint16]string
}

type Instruction interface {
	Address() uint16
	Print(printer InstructionPrinter) error
	Type() InstructionType
}

const Chip8ProgStartAddr = 0x200

type InstructionType uint8

/**
NNN: address
NN: 8-bit constant
N: 4-bit constant
X and Y: 4-bit register identifier
PC : Program Counter
I : 12bit register (For memory address) (Similar to void pointer);
VN: One of the 16 available variables. N may be 0 to F (hexadecimal);
*/

const (
	OP_0NNN InstructionType = iota // SYS addr              AddrIns
	OP_00E0                        // CLS                   ZeroIns
	OP_00EE                        // RET                   ZeroIns
	OP_1NNN                        // JP addr(NNN)          AddrIns
	OP_2NNN                        // CALL addr(NNN)        AddrIns
	OP_3XNN                        // SE Vx, 0xNN           OneRegConstIns
	OP_4XNN                        // SNE Vx, 0xNN          OneRegConstIns
	OP_5XY0                        // SE Vx, Vy             TwoRegIns
	OP_6XNN                        // LD Vx, 0xNN           OneRegConstIns
	OP_7XNN                        // ADD Vx, 0xNN          OneRegConstIns
	OP_8XY0                        // LD Vx, Vy            TwoRegIns
	OP_8XY1                        // OR, Vx, Vy           TwoRegIns
	OP_8XY2                        // AND Vx, Vy           TwoRegIns
	OP_8XY3                        // XOR Vx, Vy           TwoRegIns
	OP_8XY4                        // ADD Vx, Vy           TwoRegIns
	OP_8XY5                        // SUB Vx, Vy           TwoRegIns
	OP_8XY6                        // SHR Vx(, Vy)	        TwoRegIns
	OP_8XY7                        // SUBN Vx, Vy          TwoRegIns
	OP_8XYE                        // SHL Vx(, Vy)	        TwoRegIns
	OP_9XY0                        // SNE Vx, Vy           TwoRegIns
	OP_ANNN                        // LD I, addr(NNN)      AddrIns
	OP_BNNN                        // JP V0, addr(NNN)	    AddrIns
	OP_CXNN                        // RND Vx, 0xNN	        OneRegConstIns
	OP_DXYN                        // DRW Vx, Vy, 0xN      TwoRegConstIns
	OP_EX9E                        // SKP Vx               OneRegIns
	OP_EXA1                        // SKNP Vx              OneRegIns
	OP_FX07                        // LD Vx, DT            OneRegIns
	OP_FX0A                        // LD Vx, K	            OneRegIns
	OP_FX15                        // LD DT, Vx            OneRegIns
	OP_FX18                        // LD ST, Vx            OneRegIns
	OP_FX1E                        // ADD I, Vx            OneRegIns
	OP_FX29                        // LD F, Vx	            OneRegIns
	OP_FX33                        // LD B, Vx	            OneRegIns
	OP_FX55                        // LD [I], Vx           OneRegIns
	OP_FX65                        // LD Vx, [I]           OneRegIns
	OP_INVALID
)

func DecodeInstruction(b1 byte, b2 byte) (op InstructionType, r1 byte, r2 byte, r3 byte) {
	var r0 = b1 >> 4
	r1 = b1 & 0xf
	r2 = b2 >> 4
	r3 = b2 & 0xf
	switch r0 {
	case 0x0:
		if r1 == 0x0 && r2 == 0xE {
			if r3 == 0x0 {
				op = OP_00E0
				return
			}
			if r3 == 0xE {
				op = OP_00EE
				return
			}
		}
		op = OP_0NNN
		return
	case 0x1:
		op = OP_1NNN
		return
	case 0x2:
		op = OP_2NNN
		return
	case 0x3:
		op = OP_3XNN
		return
	case 0x4:
		op = OP_4XNN
		return
	case 0x5:
		op = OP_5XY0
		return
	case 0x6:
		op = OP_6XNN
		return
	case 0x7:
		op = OP_7XNN
		return
	case 0x8:
		switch r3 {
		case 0x0:
			op = OP_8XY0
			return
		case 0x1:
			op = OP_8XY1
			return
		case 0x2:
			op = OP_8XY2
			return
		case 0x3:
			op = OP_8XY3
			return
		case 0x4:
			op = OP_8XY4
			return
		case 0x5:
			op = OP_8XY5
			return
		case 0x6:
			op = OP_8XY6
			return
		case 0x7:
			op = OP_8XY7
			return
		case 0xE:
			op = OP_8XYE
			return
		}
	case 0x9:
		if r3 == 0x0 {
			op = OP_9XY0
			return
		}
	case 0xA:
		op = OP_ANNN
		return
	case 0xB:
		op = OP_BNNN
		return
	case 0xC:
		op = OP_CXNN
		return
	case 0xD:
		op = OP_DXYN
		return
	case 0xE:
		if r2 == 0x9 && r3 == 0xE {
			op = OP_EX9E
			return
		}
		if r2 == 0xA && r3 == 0x1 {
			op = OP_EXA1
			return
		}
	case 0xF:
		if r2 == 0x0 {
			if r3 == 0x7 {
				op = OP_FX07
				return
			}
			if r3 == 0xA {
				op = OP_FX0A
				return
			}
		} else if r2 == 0x1 {
			if r3 == 0x5 {
				op = OP_FX15
				return
			}
			if r3 == 0x8 {
				op = OP_FX18
				return
			}
			if r3 == 0xE {
				op = OP_FX1E
				return
			}
		} else if r2 == 0x2 {
			if r3 == 0x9 {
				op = OP_FX29
				return
			}
		} else if r2 == 0x3 {
			if r3 == 0x3 {
				op = OP_FX33
				return
			}
		} else if r2 == 0x5 {
			if r3 == 0x5 {
				op = OP_FX55
				return
			}
		} else if r2 == 0x6 {
			if r3 == 0x5 {
				op = OP_FX65
				return
			}
		}
	}
	op = OP_INVALID
	return
}

var InstructionTypeNames = map[InstructionType]string{
	OP_0NNN: "SYS",
	OP_00E0: "CLS",
	OP_00EE: "RET",
	OP_1NNN: "JP",
	OP_2NNN: "CALL",
	OP_3XNN: "SE",
	OP_4XNN: "SNE",
	OP_5XY0: "SE",
	OP_6XNN: "LD",
	OP_7XNN: "ADD",
	OP_8XY0: "LD",
	OP_8XY1: "OR",
	OP_8XY2: "AND",
	OP_8XY3: "XOR",
	OP_8XY4: "ADD",
	OP_8XY5: "SUB",
	OP_8XY6: "SHR",
	OP_8XY7: "SUBN",
	OP_8XYE: "SHL",
	OP_9XY0: "SNE",
	OP_ANNN: "LD",
	OP_BNNN: "JP",
	OP_CXNN: "RND",
	OP_DXYN: "DRW",
	OP_EX9E: "SKP",
	OP_EXA1: "SKNP",
	OP_FX07: "LD",
	OP_FX0A: "LD",
	OP_FX15: "LD",
	OP_FX18: "LD",
	OP_FX1E: "ADD",
	OP_FX29: "LD",
	OP_FX33: "LD",
	OP_FX55: "LD",
	OP_FX65: "LD",
}

// AddrIns follow `0nnn` form
type AddrIns struct {
	addr   uint16
	op     InstructionType
	target uint16
}

func NewAddrIns(addr uint16, op InstructionType, r1 byte, r2 byte, r3 byte) Instruction {
	return AddrIns{
		addr:   addr,
		op:     op,
		target: (uint16(r1) << 8) | (uint16(r2) << 4) | uint16(r3),
	}
}

func (a AddrIns) Type() InstructionType {
	return a.op
}

func (a AddrIns) Address() uint16 {
	return a.addr
}

func (a AddrIns) Print(printer InstructionPrinter) (err error) {
	v, ok := printer.labelMap[a.target]
	if !ok {
		v = fmt.Sprintf("@0x%03x", a.target)
	}
	_, err = fmt.Fprintf(printer.writer, "    %-4s  %s\n", InstructionTypeNames[a.op], v)
	return
}

type ZeroIns struct {
	addr uint16
	op   InstructionType
}

func NewZeroIns(addr uint16, op InstructionType, r1 byte, r2 byte, r3 byte) Instruction {
	return ZeroIns{
		addr: addr,
		op:   op,
	}
}

func (z ZeroIns) Type() InstructionType {
	return z.op
}

func (z ZeroIns) Address() uint16 {
	return z.addr
}

func (z ZeroIns) Print(printer InstructionPrinter) (err error) {
	_, err = fmt.Fprintf(printer.writer, "    %-4s\n", InstructionTypeNames[z.op])
	return
}

// OneRegIns follow `EX9E` form
type OneRegIns struct {
	addr uint16
	op   InstructionType
	reg  uint8
}

func NewOneRegIns(addr uint16, op InstructionType, r1 byte, r2 byte, r3 byte) Instruction {
	return OneRegIns{
		addr: addr,
		op:   op,
		reg:  r1,
	}
}

func (o OneRegIns) Type() InstructionType {
	return o.op
}

func (o OneRegIns) Address() uint16 {
	return o.addr
}

func (o OneRegIns) Print(printer InstructionPrinter) (err error) {
	_, err = fmt.Fprintf(printer.writer, "    %-4s  V%d\n", InstructionTypeNames[o.op], o.reg)
	return
}

// OneRegConstIns follow `3XNN` form
type OneRegConstIns struct {
	addr uint16
	op   InstructionType
	reg  uint8
	num  uint16
}

func NewOneRegConstIns(addr uint16, op InstructionType, r1 byte, r2 byte, r3 byte) Instruction {
	return OneRegConstIns{
		addr: addr,
		op:   op,
		reg:  r1,
		num:  (uint16(r2) << 4) | uint16(r3),
	}
}

func (o OneRegConstIns) Type() InstructionType {
	return o.op
}

func (o OneRegConstIns) Address() uint16 {
	return o.addr
}

func (o OneRegConstIns) Print(printer InstructionPrinter) (err error) {
	_, err = fmt.Fprintf(printer.writer, "    %-4s  V%d, 0x%02x\n", InstructionTypeNames[o.op], o.reg, o.num)
	return
}

// TwoRegIns follow `5XY0` form
type TwoRegIns struct {
	addr uint16
	op   InstructionType
	reg1 uint8
	reg2 uint8
}

func NewTwoRegIns(addr uint16, op InstructionType, r1 byte, r2 byte, r3 byte) Instruction {
	return TwoRegIns{
		addr: addr,
		op:   op,
		reg1: r1,
		reg2: r2,
	}
}

func (t TwoRegIns) Type() InstructionType {
	return t.op
}

func (t TwoRegIns) Address() uint16 {
	return t.addr
}

func (t TwoRegIns) Print(printer InstructionPrinter) (err error) {
	_, err = fmt.Fprintf(printer.writer, "    %-4s  V%d, V%d\n", InstructionTypeNames[t.op], t.reg1, t.reg2)
	return
}

// TwoRegConstIns follow `DXYN` form
type TwoRegConstIns struct {
	addr uint16
	op   InstructionType
	reg1 uint8
	reg2 uint8
	num  uint8
}

func NewTwoRegConstIns(addr uint16, op InstructionType, r1 byte, r2 byte, r3 byte) Instruction {
	return TwoRegConstIns{
		addr: addr,
		op:   op,
		reg1: r1,
		reg2: r2,
		num:  r3,
	}
}

func (t TwoRegConstIns) Type() InstructionType {
	return t.op
}

func (t TwoRegConstIns) Address() uint16 {
	return t.addr
}

func (t TwoRegConstIns) Print(printer InstructionPrinter) (err error) {
	_, err = fmt.Fprintf(printer.writer, "    %-4s  V%d, V%d, 0x%x\n", InstructionTypeNames[t.op], t.reg1, t.reg2, t.num)
	return
}

type InvalidIns struct {
	addr uint16
	b1   byte
	b2   byte
}

func (i InvalidIns) Address() uint16 {
	return i.addr
}

func (i InvalidIns) Print(printer InstructionPrinter) (err error) {
	_, err = fmt.Fprintf(printer.writer, "    0x%02x%02x\n", i.b1, i.b2)
	return
}

func (i InvalidIns) Type() InstructionType {
	return OP_INVALID
}
