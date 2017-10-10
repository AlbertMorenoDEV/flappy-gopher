package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"fmt"
	"sync"
)

const (
	gravity = 0.1
	jumpSpeed = 4
)

type bird struct {
	mu sync.RWMutex

	time int
	textures []*sdl.Texture

	x, y int32
	w, h int32
	speed float64
	dead bool
}

func newBird(render *sdl.Renderer, screen *screen) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("resources/images/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(render, path)
		if err != nil {
			return nil, fmt.Errorf("could not load texture image: %v", err)
		}
		textures = append(textures, texture)
	}
	return &bird{textures: textures, x: 10, y: screen.h/2, w: 50, h: 43}, nil
}

func (bird *bird) update() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.time++
	bird.y -= int32(bird.speed)
	if bird.y < 0 {
		bird.dead = true
	}
	bird.speed += gravity
}

func (bird *bird) paint(render *sdl.Renderer, screen *screen) error {
	bird.mu.RLock()
	defer bird.mu.RUnlock()

	rect := &sdl.Rect{X: bird.x, Y: (screen.h - bird.y) - bird.h/2, W: bird.w, H: bird.h}

	i := bird.time / 10 % len(bird.textures)
	if err := render.Copy(bird.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}
	return nil
}

func (bird *bird) restart(screen *screen) {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.y = screen.h/2
	bird.speed = 0
	bird.dead = false
}

func (bird *bird) destroy() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	for _, texture := range bird.textures {{
		texture.Destroy()
	}}
}

func (bird *bird) isDead() bool {
	bird.mu.RLock()
	defer bird.mu.RUnlock()

	return bird.dead
}

func (bird *bird) jump() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.speed = -jumpSpeed
}

func (bird *bird) touch(pipe *pipe, screen *screen) {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	if pipe.x > bird.x + bird.w { // too far right
		return
	}

	if pipe.x + pipe.w < bird.x { // too far left
		return
	}

	if !pipe.inverted && pipe.h < bird.y-bird.h/2 { // pipe is too low
		return
	}

	if pipe.inverted && screen.h-pipe.h > bird.y+bird.h/2 { // pipe is too high
		return
	}

	bird.dead = true
}