package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"context"
	"time"
)

type scene struct {
	time int
	background *sdl.Texture
	birds []*sdl.Texture
}

func newScene(render *sdl.Renderer) (*scene, error) {
	texture, err := img.LoadTexture(render, "resources/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	var birds []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("resources/images/bird_frame_%d.png", i)
		bird, err := img.LoadTexture(render, path)
		if err != nil {
			return nil, fmt.Errorf("could not load bird image: %v", err)
		}
		birds = append(birds, bird)
	}

	return &scene{background: texture, birds: birds}, nil
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
	scene.time++

	render.Clear()

	if err := render.Copy(scene.background, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	rect := &sdl.Rect{X:10, Y:300-43/2, W:50, H:43}

	i := scene.time / 10 % len(scene.birds)
	if err := render.Copy(scene.birds[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}

	render.Present()

	return nil
}

func (scene *scene) destroy() {
	scene.background.Destroy()
}