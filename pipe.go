package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	mu       sync.RWMutex
	posX     int32
	posY     int32
	height   int32
	width    int32
	speed    int32
	inverted bool
	ren      *sdl.Renderer
	pTexture *sdl.Texture
}

const (
	constPipeWidth  = 50
	constPipeHeight = 300
)

func newPipe(ren *sdl.Renderer) (*pipe, error) {
	pTexture, err := img.LoadTexture(ren, "res/images/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("cannot load pipe: %v", err)
	}
	return &pipe{
		posX:     ConstScreenWidth - constPipeWidth,
		posY:     ConstScreenHeight - constPipeHeight,
		height:   constPipeHeight,
		width:    constPipeWidth,
		speed:    10,
		inverted: false,
		pTexture: pTexture,
		ren:      ren,
	}, nil
}

func (p *pipe) paint() error {
	p.mu.RLock()
	rect := &sdl.Rect{
		X: p.posX,
		Y: p.posY,
		W: p.width,
		H: p.height,
	}

	if err := p.ren.Copy(p.pTexture, nil, rect); err != nil {
		return fmt.Errorf("cannot copy pipe into renderer: %v", err)
	}
	p.mu.RUnlock()
	return nil
}

func (p *pipe) update() {
	p.mu.Lock()
	p.posX -= p.speed
	p.mu.Unlock()
}

func (p *pipe) restart() {
	p.mu.Lock()
	p.posX = ConstScreenWidth - constPipeWidth
	p.mu.Unlock()
}

func (p *pipe) destroy() {
	p.mu.Lock()
	p.pTexture.Destroy()
	p.mu.Unlock()
}
