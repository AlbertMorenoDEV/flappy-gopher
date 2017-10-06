package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"github.com/veandco/go-sdl2/img"
)

type bird struct {
	time int
	textures []*sdl.Texture
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
	return &bird{textures: textures}, nil
}

func (bird *bird) paint(render *sdl.Renderer) error {
	bird.time++

	rect := &sdl.Rect{X:10, Y:300-43/2, W:50, H:43}

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