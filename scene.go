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
	pipes *pipes
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

	pipes, err := newPipes(render)
	if err != nil {
		return nil, err
	}

	return &scene{background: texture, bird: bird, pipes: pipes}, nil
}

func (scene *scene) run(events <-chan sdl.Event, render *sdl.Renderer) chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
				case e := <-events:
					if done := scene.handleEvent(e); done {
						return
					}
				case <-tick:
					scene.update()

					if scene.bird.isDead() {
						drawTitle(render, "Game Over")
						time.Sleep(time.Second)
						scene.restart()
					}

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

func (scene *scene) update() {
	scene.bird.update()
	scene.pipes.update()
	scene.pipes.touch(scene.bird)
}

func (scene *scene) restart() {
	scene.bird.restart()
	scene.pipes.restart()
}

func (scene *scene) paint(render *sdl.Renderer) error {
	render.Clear()

	if err := render.Copy(scene.background, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := scene.bird.paint(render); err != nil {
		return err
	}

	if err := scene.pipes.paint(render); err != nil {
		return err
	}

	render.Present()

	return nil
}

func (scene *scene) destroy() {
	scene.background.Destroy()
	scene.bird.destroy()
	scene.pipes.destroy()
}