package main

import "github.com/veandco/go-sdl2/sdl"

type SDLDevice struct {
}

func (sdlDevice *SDLDevice) Setup() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}
	return nil
}

func (sdlDevice *SDLDevice) Teardown() error {
	sdl.Quit()
	return nil
}
