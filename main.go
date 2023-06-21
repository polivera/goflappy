package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	var err error
	//	Init SDL
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("cannot initialize SDL: %v", err)
	}
	defer sdl.Quit()
	// Init fonts
	if err = ttf.Init(); err != nil {
		return fmt.Errorf("cannot initialize Font: %v", err)
	}

	win, ren, err := sdl.CreateWindowAndRenderer(
		ConstScreenWidth,
		ConstScreenHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return fmt.Errorf("cannot create window renderer: %v", err)
	}
	defer win.Destroy()

	scene, err := newScene(ren)
	if err != nil {
		return fmt.Errorf("cannot create a scene")
	}

	events := make(chan sdl.Event)
	errc := scene.run(events, ren)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

func drawBackground(ren *sdl.Renderer) error {
	ren.Clear()
	defer ren.Present()

	bgTexture, err := img.LoadTexture(ren, "res/images/background.png")
	if err != nil {
		return fmt.Errorf("cannot load background: %v", err)
	}
	defer bgTexture.Destroy()

	if err = ren.Copy(bgTexture, nil, nil); err != nil {
		return fmt.Errorf("cannot copy background to renderer: %v", err)
	}

	return nil
}

func drawTitle(ren *sdl.Renderer) error {
	ren.Clear()
	defer ren.Present()

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

	texture, err := ren.CreateTextureFromSurface(sfc)
	if err != nil {
		return fmt.Errorf("cannot create a texture from surface for font: %v", err)
	}
	defer texture.Destroy()

	err = ren.Copy(texture, nil, nil)
	if err != nil {
		return fmt.Errorf("cannot copy font texture to renderer: %v", err)
	}
	return nil
}
