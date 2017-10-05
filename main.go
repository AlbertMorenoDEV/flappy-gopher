package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
	"time"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/img"
)

func main() {
	if err := run() ; err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}

	window, render, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer window.Destroy()

	if err := drawTitle(render); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(5 * time.Second)

	if err := drawBackground(render); err != nil {
		return fmt.Errorf("could now draw background: %v", err)
	}

	time.Sleep(5 * time.Second)

	return nil
}

func drawBackground(render *sdl.Renderer) error {
	render.Clear()

	texture, err := img.LoadTexture(render, "resources/images/background.png")
	if err != nil {
		return fmt.Errorf("could not load background image: %v", err)
	}
	defer texture.Destroy()

	if err := render.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	render.Present()

	return nil
}

func drawTitle(render *sdl.Renderer) error {
	render.Clear()

	font, err := ttf.OpenFont("resources/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer font.Close()

	color := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	surface, err := font.RenderUTF8_Solid("Flappy Gopher", color)
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
