package main

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/kong"
	"os"
)

type CLIRun struct {
	Path string `arg:"positional" required:"" help:"Path to CHIP-8 ROM"`
}

type CLIDisasm struct {
	Path string `arg:"positional" required:"" help:"Path to CHIP-8 ROM"`
}

var CLI struct {
	Run CLIRun `cmd:"" help:"Run CHIP-8 ROM"`

	Disasm CLIDisasm `cmd:"" help:"Disassemble CHIP-8 ROM"`
}

func (r *CLIRun) Run() error {
	device := SDLDevice{}
	err := device.Setup()
	if err != nil {
		return fmt.Errorf("device setup error: %v\n", err)
	}
	defer device.Teardown()

	buf, err := os.ReadFile(r.Path)
	if err != nil {
		return fmt.Errorf("run error: %v\n", err)
	}
	reader := bytes.NewReader(buf)
	vm, err := NewChip8VM(reader, &device)
	if err != nil {
		return fmt.Errorf("run error: %v\n", err)
	}
	vm.Dump(os.Stdout)
	if err = vm.Run(); err != nil {
		return fmt.Errorf("run error: %v\n", err)
	}
	return nil
}

func (d *CLIDisasm) Run() error {
	buf, err := os.ReadFile(d.Path)
	if err != nil {
		return fmt.Errorf("disasm error: %v\n", err)
	}
	reader := bytes.NewReader(buf)
	err = Disassemble(reader, os.Stdout)
	if err != nil {
		return fmt.Errorf("disasm error: %v\n", err)
	}
	return nil
}

func main() {
	ctx := kong.Parse(&CLI, kong.UsageOnError())
	err := ctx.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
