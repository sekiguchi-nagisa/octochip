package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type Display interface {
	Clear()
	Draw(x byte, y byte, addr uint16, n byte) uint8
}

type Keyboard interface {
	KeyPressed(key uint8) bool
	WaitKey() uint8
	Poll()
}

const Chip8RAMSize = 4096

var fontSpriteSet = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8VM struct {
	ram      [Chip8RAMSize]uint8
	reg      [16]uint8
	I        uint16     // for memory address
	dt       uint8      // for delay timer
	st       uint8      // for sound timer
	pc       uint16     // program counter
	sp       uint8      // stack pointer
	stack    [16]uint16 // maintains return address
	rng      *rand.Rand
	display  Display
	keyboard Keyboard
}

func NewChip8VM(reader io.Reader, display Display, keyboard Keyboard) (*Chip8VM, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	vm := Chip8VM{}
	vm.rng = rand.New(rand.NewSource(42))
	vm.display = display
	vm.keyboard = keyboard

	// full preset font sprite
	for i := 0; i < len(fontSpriteSet); i++ {
		vm.ram[i] = fontSpriteSet[i]
	}

	vm.pc = Chip8ProgStartAddr
	for i := 0; i < len(buf); i++ {
		vm.ram[Chip8ProgStartAddr+i] = buf[i]
	}
	return &vm, nil
}

// Dump dump internal state
func (vm *Chip8VM) Dump(writer io.Writer) {
	// dump registers
	for i, u := range vm.reg {
		_, _ = fmt.Fprintf(writer, "V%x=0x%02X", i, u)
		if i+1%8 == 0 {
			_, _ = fmt.Fprint(writer, "\n")
		} else {
			_, _ = fmt.Fprint(writer, " ")
		}
	}
	_, _ = fmt.Fprintf(writer, "I=0x%03X, pc=0x%04X, sp=0x%02X\n", vm.I, vm.pc, vm.sp)
	_, _ = fmt.Fprint(writer, "stack[")
	for i, v := range vm.stack {
		if i > 0 {
			_, _ = fmt.Fprint(writer, ", ")
		}
		_, _ = fmt.Fprintf(writer, "0x%04X", v)
	}
	_, _ = fmt.Fprint(writer, "]\n")
	_, _ = fmt.Fprintf(writer, "DT=%d, ST=%d\n", vm.dt, vm.st)
}

// Run entry point
func (vm *Chip8VM) Run() {
	prevMicroSec := time.Now().UnixMicro()
	for int(vm.pc) <= len(vm.ram) {
		vm.dispatchSingleIns()
		vm.keyboard.Poll()

		// decrement delay/sound timer
		curMicroSec := time.Now().UnixMicro()
		dec := (curMicroSec - prevMicroSec) / 16667
		if dec <= int64(vm.dt) {
			vm.dt = uint8(dec)
		} else {
			vm.dt = 0
		}
		if dec <= int64(vm.st) {
			vm.st = uint8(dec)
		} else {
			vm.st = 0
		}
		prevMicroSec = curMicroSec
	}
}

func (vm *Chip8VM) dispatchSingleIns() {
	if vm.pc >= Chip8RAMSize {
		panic("program counter overflow: " + strconv.Itoa(int(vm.pc)))
	}
	if vm.sp > 16 {
		panic("stack pointer overflow: " + strconv.Itoa(int(vm.sp)))
	}

	b1 := vm.ram[vm.pc]
	b2 := vm.ram[vm.pc+1]
	op, r1, r2, r3 := DecodeInstruction(b1, b2)
	vm.pc += 2

	targetAddr := (uint16(r1) << 8) | (uint16(r2) << 4) | uint16(r3)
	num := (r2 << 4) | r3
	switch op {
	case OP_0NNN:
		// do nothing
	case OP_00E0:
		vm.display.Clear()
	case OP_00EE:
		vm.pc = vm.stack[vm.sp-1]
		vm.sp--
	case OP_1NNN:
		vm.pc = targetAddr
	case OP_2NNN:
		vm.sp++
		vm.stack[vm.sp-1] = vm.pc
		vm.pc = targetAddr
	case OP_3XNN:
		if vm.reg[r1] == num {
			vm.pc += 2
		}
	case OP_4XNN:
		if vm.reg[r1] != num {
			vm.pc += 2
		}
	case OP_5XY0:
		if vm.reg[r1] == vm.reg[r2] {
			vm.pc += 2
		}
	case OP_6XNN:
		vm.reg[r1] = num
	case OP_7XNN:
		vm.reg[r1] += num
	case OP_8XY0:
		vm.reg[r1] = vm.reg[r2]
	case OP_8XY1:
		vm.reg[r1] |= vm.reg[r2]
	case OP_8XY2:
		vm.reg[r1] &= vm.reg[r2]
	case OP_8XY3:
		vm.reg[r1] ^= vm.reg[r2]
	case OP_8XY4:
		if vm.reg[r1] > uint8(math.MaxUint8)-vm.reg[r1] { // overflow
			vm.reg[0xF] = 1
		} else {
			vm.reg[0xF] = 0
		}
		vm.reg[r1] += vm.reg[r2]
	case OP_8XY5:
		if vm.reg[r1] < vm.reg[r2] { // underflow
			vm.reg[0xF] = 0
		} else {
			vm.reg[0xF] = 1
		}
		vm.reg[r1] -= vm.reg[r2]
	case OP_8XY6:
		vm.reg[0xF] = vm.reg[r1] & 0x1
		vm.reg[r1] >>= 1
	case OP_8XY7:
		if vm.reg[r1] > vm.reg[r2] { // underflow
			vm.reg[0xF] = 0
		} else {
			vm.reg[0xF] = 1
		}
		vm.reg[r1] = vm.reg[r2] - vm.reg[r1]
	case OP_8XYE:
		vm.reg[0xF] = vm.reg[r1] >> 7
		vm.reg[r1] <<= 1
	case OP_9XY0:
		if vm.reg[r1] != vm.reg[r2] {
			vm.pc += 2
		}
	case OP_ANNN:
		vm.I = targetAddr
	case OP_BNNN:
		vm.pc = uint16(vm.reg[0]) + targetAddr
	case OP_CXNN:
		vm.reg[r1] = uint8(vm.rng.Intn(256)) & num
	case OP_DXYN:
		vm.reg[0xF] = vm.display.Draw(vm.reg[r1], vm.reg[r2], vm.I, r3)
	case OP_EX9E:
		if vm.keyboard.KeyPressed(vm.reg[r1]) {
			vm.pc += 2
		}
	case OP_EXA1:
		if !vm.keyboard.KeyPressed(vm.reg[r1]) {
			vm.pc += 2
		}
	case OP_FX07:
		vm.reg[r1] = vm.dt
	case OP_FX0A:
		vm.reg[r1] = vm.keyboard.WaitKey()
	case OP_FX15:
		vm.dt = vm.reg[r1]
	case OP_FX18:
		vm.st = vm.reg[r1]
	case OP_FX1E:
		vm.I += uint16(vm.reg[r1])
	case OP_FX29:
		vm.I = uint16(vm.reg[r1]) * 5
	case OP_FX33:
		vm.ram[vm.I] = vm.reg[r1] / 100
		vm.ram[vm.I+1] = vm.reg[r1] / 10 % 10
		vm.ram[vm.I+2] = vm.reg[r1] % 100 % 10
	case OP_FX55:
		for i := 0; i <= int(r1); i++ {
			vm.ram[vm.I+uint16(i)] = vm.reg[i]
		}
	case OP_FX65:
		for i := 0; i <= int(r1); i++ {
			vm.reg[i] = vm.ram[vm.I+uint16(i)]
		}
	default:
		ss := fmt.Sprintf("invalid opcode: %d, 0x%02x%02x", op, b1, b2)
		panic(ss)
	}
}