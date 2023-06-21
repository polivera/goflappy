package main

import (
	"fmt"
	"os"
	"runtime"

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

