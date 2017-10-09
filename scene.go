package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"time"
	"log"
)

type scene struct {
	background *sdl.Texture
	bird *bird
}

func newScene(render *sdl.Renderer) (*scene, error) {
	texture, err := img.LoadTexture(render, "resources/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	bird, err := newBird(render)
	if err != nil {
		return nil, err
	}


	return &scene{background: texture, bird: bird}, nil
}

func (scene *scene) run(events <-chan sdl.Event, render *sdl.Renderer) chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Microsecond)
		for {
			select {
				case e := <-events:
					if done := scene.handleEvent(e); done {
						return
					}
				case <-tick:
					if err := scene.paint(render); err != nil {
						errc <- err
					}
			}
		}
	}()

	return errc
}

func (scene *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.MouseButtonEvent:
			scene.bird.jump()
		case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent:
		default:
			log.Printf("unkown event %T", event)
	}
	return false
}

func (scene *scene) paint(render *sdl.Renderer) error {
	render.Clear()

	if err := render.Copy(scene.background, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := scene.bird.paint(render); err != nil {
		return err
	}

	render.Present()

	return nil
}

func (scene *scene) destroy() {
	scene.background.Destroy()
	scene.bird.destroy()
}