package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
	"time"
	"github.com/veandco/go-sdl2/ttf"
	"runtime"
)

const (
	w = 800
	h = 600
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

	screen := newScreen(w, h)

	window, render, err := sdl.CreateWindowAndRenderer(int(screen.w), int(screen.h), sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer window.Destroy()

	if err := screen.drawTitle(render, "Flappy Gopher"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(2 * time.Second)

	scene, err := newScene(render, screen)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer scene.destroy()

	events := make(chan sdl.Event)
	errc := scene.run(events, render, screen)

	runtime.LockOSThread()
	for {
		select {
			case events <- sdl.WaitEvent():
			case err := <-errc:
				return err
		}
	}
}
