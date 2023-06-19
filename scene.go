package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	tm    uint8
	bg    *sdl.Texture
	birds *bird
}

// newScene
func newScene(ren *sdl.Renderer) (*scene, error) {
	bgTexture, err := img.LoadTexture(ren, "res/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("cannot load background: %v", err)
	}

	birds, err := newBird(ren)
	if err != nil {
		return nil, fmt.Errorf("cannot load bird textures: %v:", err)
	}
	

	return &scene{
		bg:    bgTexture,
		birds: birds,
	}, nil
}

// run
func (s *scene) run(ctx context.Context, ren *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(100 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(ren); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

// paint
func (s *scene) paint(ren *sdl.Renderer) error {
	ren.Clear()

	if err := ren.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}

	if err := s.birds.paint(); err != nil {
		return fmt.Errorf("could not create bird on scene: %v", err)
	}

	ren.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
}
