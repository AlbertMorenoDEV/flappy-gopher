package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
	"time"
	"github.com/veandco/go-sdl2/ttf"
	"runtime"
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

	if err := drawTitle(render, "Flappy Gopher"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(2 * time.Second)

	scene, err := newScene(render)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer scene.destroy()

	events := make(chan sdl.Event)
	errc := scene.run(events, render)

	runtime.LockOSThread()
	for {
		select {
			case events <- sdl.WaitEvent():
			case err := <-errc:
				return err
		}
	}
}

func drawTitle(render *sdl.Renderer, text string) error {
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
