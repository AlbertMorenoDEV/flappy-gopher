package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"sync"
)

type pipe struct {
	mu sync.RWMutex

	texture *sdl.Texture

	x int32
	h int32
	w int32
	speed int32
	inverted bool
}

func newPipe(render *sdl.Renderer) (*pipe, error) {
	texture, err := img.LoadTexture(render, "resources/images/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}

	return &pipe {
		texture: texture,
		x: 400,
		h: 200,
		w: 50,
		speed: 1,
		inverted: true,
	}, nil
}

func (pipe *pipe) paint(render *sdl.Renderer) error {
	pipe.mu.RLock()
	defer pipe.mu.RUnlock()

	rect := &sdl.Rect{X: pipe.x, Y: 600 - pipe.h, W: pipe.w, H: pipe.h}
	flip := sdl.FLIP_NONE
	if pipe.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := render.CopyEx(pipe.texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}

	return nil
}

func (pipe *pipe) restart() {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()

	pipe.x = 400
}

func (pipe *pipe) update() {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()

	pipe.x -= pipe.speed
}

func (pipe *pipe) destroy() {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()

	pipe.texture.Destroy()
}