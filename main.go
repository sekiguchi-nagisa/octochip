package main

import (
	"fmt"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Run struct {
		Path string `arg:"positional" required:"" help:"Path to CHIP-8 ROM"`
	} `cmd:"" help:"Run CHIP-8 ROM"`

	Disasm struct {
		Path string `arg:"positional" required:"" help:"Path to CHIP-8 ROM"`
	} `cmd:"" help:"Disassemble CHIP-8 ROM"`

	Build struct {
		Path string `arg:"positional" required:"" help:"Path to CHIP-8 Program"`
	} `cmd:"" help:"Build CHIP-8 ROM"`
}

func main() {
	ctx := kong.Parse(&CLI, kong.UsageOnError())
	switch ctx.Command() {
	case "run <path>":
		fmt.Println("FIXME: run command")
	case "disasm <path>":
		fmt.Println("FIXME: disasm command")
	case "build <path>":
		fmt.Println("FIXME: build command")
	default:
		panic(ctx.Command())
	}
}
