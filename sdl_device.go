package main

import "github.com/veandco/go-sdl2/sdl"

type SDLDevice struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
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
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, 64, 32)
	if err != nil {
		return err
	}
	sdlDevice.window = window
	sdlDevice.renderer = renderer
	sdlDevice.texture = texture
	return nil
}

func (sdlDevice *SDLDevice) Teardown() error {
	err := sdlDevice.window.Destroy()
	if err != nil {
		return err
	}
	err = sdlDevice.renderer.Destroy()
	if err != nil {
		return err
	}
	err = sdlDevice.texture.Destroy()
	if err != nil {
		return err
	}
	sdl.Quit()
	return nil
}
