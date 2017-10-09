package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"github.com/veandco/go-sdl2/img"
)

const (
	gravity = 0.25
	jumpSpeed = 5
)

type bird struct {
	time int
	textures []*sdl.Texture

	y, speed float64
}

func newBird(render *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("resources/images/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(render, path)
		if err != nil {
			return nil, fmt.Errorf("could not load texture image: %v", err)
		}
		textures = append(textures, texture)
	}
	return &bird{textures: textures, y: 300}, nil
}

func (bird *bird) paint(render *sdl.Renderer) error {
	bird.time++
	bird.y -= bird.speed
	if bird.y < 0 {
		bird.speed = -bird.speed
		bird.y = 0
	}
	bird.speed += gravity

	rect := &sdl.Rect{X:10, Y: (600 - int32(bird.y)) - 43/2, W:50, H:43}

	i := bird.time / 10 % len(bird.textures)
	if err := render.Copy(bird.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}
	return nil
}

func (bird *bird) destroy() {
	for _, texture := range bird.textures {{
		texture.Destroy()
	}}
}

func (bird *bird) jump() {
	bird.speed = -jumpSpeed
}