package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const gravity = 3.3

type bird struct {
	time     uint8
	textures []*sdl.Texture
	ren      *sdl.Renderer
	y, speed float64
}

// newBird
func newBird(ren *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/images/bird_frame_%d.png", i)

		texture, err := img.LoadTexture(ren, path)
		if err != nil {
			return nil, fmt.Errorf("cannot load bird frame %d: %v", i, err)
		}
		textures = append(textures, texture)
	}
	return &bird{textures: textures, ren: ren, y: ConstBirdStartingPoint}, nil
}

func (br *bird) paint() error {
	br.time++
	br.y -= br.speed
	br.speed += gravity
	if br.y < -ConstScreenHeight {
		br.speed *= -1
		br.y = -ConstScreenHeight + 43
	}

	if br.time > 255 {
		br.time = 0
	}

	rect := &sdl.Rect{
		X: 20,
		Y: ((ConstScreenHeight - int32(br.y)) - 43) / 2,
		W: 50,
		H: 43,
	}

	ind := br.time % uint8(len(br.textures))
	if err := br.ren.Copy(br.textures[ind], nil, rect); err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	return nil
}
