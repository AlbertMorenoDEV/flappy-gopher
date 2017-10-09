package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"sync"
	"time"
	"math/rand"
)

type pipes struct {
	mu sync.RWMutex

	texture *sdl.Texture
	speed int32

	pipes []*pipe
}

func newPipes(render *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(render, "resources/images/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}

	pipes := &pipes {
		texture: texture,
		speed: 2,
	}

	go func() {
		for {
			pipes.mu.Lock()
			pipes.pipes = append(pipes.pipes, newPipe())
			pipes.mu.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	return pipes, nil
}

func (pipes *pipes) paint(render *sdl.Renderer) error {
	pipes.mu.RLock()
	defer pipes.mu.RUnlock()

	for _, pipe := range pipes.pipes {
		if err := pipe.paint(render, pipes.texture); err != nil {
			return err
		}
	}

	return nil
}

func (pipes *pipes) touch(bird *bird) {
	pipes.mu.RLock()
	defer pipes.mu.RUnlock()

	for _, pipe := range pipes.pipes {
		pipe.touch(bird)
	}
}

func (pipes *pipes) restart() {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	pipes.pipes = nil
}

func (pipes *pipes) update() {
	pipes.mu.RLock()
	defer pipes.mu.RUnlock()

	var rem []*pipe

	for _, pipe := range pipes.pipes {
		pipe.mu.RLock()
		pipe.x -= pipes.speed
		pipe.mu.RUnlock()
		if pipe.x+pipe.w > 0 {
			rem = append(rem, pipe)
		}
	}
}

func (pipes *pipes) destroy() {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	pipes.texture.Destroy()
}

type pipe struct {
	mu sync.RWMutex

	x int32
	h int32
	w int32
	inverted bool
}

func newPipe() (*pipe) {
	return &pipe {
		x: 800,
		h: 100 + int32(rand.Intn(300)),
		w: 50,
		inverted: rand.Float32() > 0.5,
	}
}

func (pipe *pipe) touch(bird *bird) {
	pipe.mu.RLock()
	defer pipe.mu.RUnlock()

	bird.touch(pipe)
}

func (pipe *pipe) paint(render *sdl.Renderer, texture *sdl.Texture) error {
	pipe.mu.RLock()
	defer pipe.mu.RUnlock()

	rect := &sdl.Rect{X: pipe.x, Y: 600 - pipe.h, W: pipe.w, H: pipe.h}
	flip := sdl.FLIP_NONE
	if pipe.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := render.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}

	return nil
}