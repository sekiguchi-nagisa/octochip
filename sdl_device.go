package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

/*
original keycode
1 2 3 C
4 5 6 D
7 8 9 E
A 0 B F

physical keycode
1 2 3 4
q w e r
a s d f
z x c v
*/
var keyMap = map[sdl.Keycode]uint8{
	sdl.K_1: 0x01,
	sdl.K_2: 0x02,
	sdl.K_3: 0x03,
	sdl.K_4: 0x0C,
	sdl.K_q: 0x04,
	sdl.K_w: 0x05,
	sdl.K_e: 0x06,
	sdl.K_r: 0x0d,
	sdl.K_a: 0x07,
	sdl.K_s: 0x08,
	sdl.K_d: 0x09,
	sdl.K_f: 0x0E,
	sdl.K_z: 0x0A,
	sdl.K_x: 0x00,
	sdl.K_c: 0x0B,
	sdl.K_v: 0x0F,
}

type SDLDevice struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

const scale = 8

func (sdlDevice *SDLDevice) Draw(screen *Screen) error {
	if err := sdlDevice.renderer.SetDrawColor(0, 0, 0, 1); err != nil {
		return err
	}
	if err := sdlDevice.renderer.Clear(); err != nil {
		return err
	}
	if err := sdlDevice.renderer.SetDrawColor(255, 255, 255, 1); err != nil {
		return err
	}
	for height := 0; height < ScreenHeight; height++ {
		for width := 0; width < ScreenWidth; width++ {
			pixel := screen[ScreenWidth*height+width]
			if pixel != 0 {
				err := sdlDevice.renderer.FillRect(&sdl.Rect{
					X: int32(width) * scale,
					Y: int32(height) * scale,
					W: scale,
					H: scale,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	sdlDevice.renderer.Present()
	return nil
}

func (sdlDevice *SDLDevice) PollKey(key *KeyPressed) bool {
	switch event := sdl.PollEvent().(type) {
	case *sdl.QuitEvent:
		fmt.Println("Quit")
		return false
	case *sdl.KeyboardEvent:
		switch event.Type {
		case sdl.KEYDOWN:
			if keycode, ok := keyMap[event.Keysym.Sym]; ok {
				fmt.Printf("keydown: %s => %x\n", sdl.GetKeyName(event.Keysym.Sym), keycode)
				*key = *key | (1 << keycode)
			}
			if event.Keysym.Sym == sdl.K_ESCAPE {
				fmt.Printf("Quit: %s\n", sdl.GetKeyName(event.Keysym.Sym))
				return false
			}
		case sdl.KEYUP:
			if keycode, ok := keyMap[event.Keysym.Sym]; ok {
				fmt.Printf("keyup: %s => %x\n", sdl.GetKeyName(event.Keysym.Sym), keycode)
				*key = *key & ^(1 << keycode)
			}
		}
	}
	return true
}

func (sdlDevice *SDLDevice) Setup() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}
	window, renderer, err := sdl.CreateWindowAndRenderer(640, 480, sdl.WINDOW_OPENGL)
	if err != nil {
		return err
	}
	window.SetTitle("octochip")
	sdlDevice.window = window
	sdlDevice.renderer = renderer
	return nil
}

func (sdlDevice *SDLDevice) Teardown() {
	sdl.Quit()
}
