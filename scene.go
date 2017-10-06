package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"context"
	"time"
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

func (scene *scene) run(ctx context.Context, render *sdl.Renderer) chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Microsecond){
			select {
				case <-ctx.Done():
					return
				default:
					if err := scene.paint(render); err != nil {
						errc <- err
					}
			}
		}
	}()

	return errc
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