package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type scene struct {
	tm    uint8
	bg    *sdl.Texture
	ren   *sdl.Renderer
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
		ren:   ren,
		birds: birds,
	}, nil
}

// run
func (s *scene) run(events <-chan sdl.Event, ren *sdl.Renderer) <-chan error {
	errc := make(chan error)

	if err := s.drawTitle(); err != nil {
		errc <- err
	} else {
		time.Sleep(2 * time.Second)
	}

	go func() {
		defer close(errc)
		ticker := time.Tick(50 * time.Millisecond)
		done := false
		for !done {
			select {
			case e := <-events:
				done = s.handleEvent(e)
			case <-ticker:
				if err := s.paint(ren); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) handleEvent(ev sdl.Event) bool {
	switch ev.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.birds.jump()
	default:
		fmt.Printf("unkown event %T\n", ev)
	}
	return false
}

func (s *scene) drawTitle() error {
	s.ren.Clear()
	defer s.ren.Present()

	fnt, err := ttf.OpenFont("res/fonts/playball.ttf", 720)
	if err != nil {
		return fmt.Errorf("cannot load font: %v", err)
	}
	defer fnt.Close()

	color := sdl.Color{
		R: 0x33,
		G: 0x66,
		B: 0x99,
	}
	sfc, err := fnt.RenderUTF8Solid("Flappy Gopher", color)
	if err != nil {
		return fmt.Errorf("cannot get surface for the font:  %v", err)
	}
	defer sfc.Free()

	texture, err := s.ren.CreateTextureFromSurface(sfc)
	if err != nil {
		return fmt.Errorf("cannot create a texture from surface for font: %v", err)
	}
	defer texture.Destroy()

	err = s.ren.Copy(texture, nil, nil)
	if err != nil {
		return fmt.Errorf("cannot copy font texture to renderer: %v", err)
	}
	return nil
}

// paint
func (s *scene) paint(ren *sdl.Renderer) error {
	ren.Clear()

	if err := ren.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not create background: %v", err)
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
