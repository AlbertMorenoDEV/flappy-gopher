package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type screen struct {
	w, h int32
}

func newScreen(w int32, h int32) *screen {
	return &screen{w: w, h: h}
}

func (screen *screen) drawTitle(render *sdl.Renderer, text string) error {
	render.Clear()

	font, err := ttf.OpenFont("resources/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer font.Close()

	color := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	surface, err := font.RenderUTF8_Solid(text, color)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer surface.Free()

	texture, err := render.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer texture.Destroy()

	if err := render.Copy(texture, nil, nil); err != nil {
		fmt.Errorf("could not copy texture: %v", err)
	}

	render.Present()

	return nil
}