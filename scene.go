package main

import (
	"fmt"
	"log"
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
func (s *scene) run(events <-chan sdl.Event, ren *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		ticker := time.Tick(50 * time.Millisecond)
		for {
			select {
			case e := <-events:
				log.Printf("event %T", e)
			case <-ticker:
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
