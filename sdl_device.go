package main

import "github.com/veandco/go-sdl2/sdl"

type SDLDevice struct {
}

func NewSDLDevice() *SDLDevice {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	return &SDLDevice{}
}
