package main

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/kong"
	"os"
)

var CLI struct {
	Run struct {
		Path string `arg:"positional" required:"" help:"Path to CHIP-8 ROM"`
	} `cmd:"" help:"Run CHIP-8 ROM"`

	Disasm struct {
		Path string `arg:"positional" required:"" help:"Path to CHIP-8 ROM"`
	} `cmd:"" help:"Disassemble CHIP-8 ROM"`
}

func main() {
	ctx := kong.Parse(&CLI, kong.UsageOnError())
	switch ctx.Command() {
	case "run <path>":
		device := SDLDevice{}
		err := device.Setup()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "device setup error: %v\n", err)
			os.Exit(1)
		}
		defer device.Teardown()

		buf, err := os.ReadFile(CLI.Run.Path)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "run error: %v\n", err)
			os.Exit(1)
		}
		reader := bytes.NewReader(buf)
		vm, err := NewChip8VM(reader, &device)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "run error: %v\n", err)
			os.Exit(1)
		}
		vm.Dump(os.Stdout)
		if err = vm.Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "run error: %v\n", err)
			os.Exit(1)
		}
	case "disasm <path>":
		buf, err := os.ReadFile(CLI.Disasm.Path)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "disasm error: %v\n", err)
			os.Exit(1)
		}
		reader := bytes.NewReader(buf)
		err = Disassemble(reader, os.Stdout)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "disasm error: %v\n", err)
			os.Exit(1)
		}
	default:
		panic(ctx.Command())
	}
}
