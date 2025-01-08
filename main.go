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
		fmt.Println("FIXME: run command")
	case "disasm <path>":
		//Dis
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
