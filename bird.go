package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	constSpeed   = 30
	constGravity = 3.3
)

// Patch to handle the fact that mouse down and mouse up
// are the same event
var isFirstJump = true

type bird struct {
	mu       sync.RWMutex
	dead     bool
	time     uint8
	textures []*sdl.Texture
	ren      *sdl.Renderer
	speed    float64
	posY     int32
	posX     int32
	width    int32
	height   int32
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
	return &bird{
		textures: textures,
		ren:      ren,
		posX:     20,
		posY:     ConstBirdStartingPoint,
		width:    50,
		height:   43,
	}, nil
}

func (br *bird) update() {
	br.mu.Lock()
	br.time++
	br.posY -= int32(br.speed)
	br.speed += constGravity
	if br.posY < -ConstScreenHeight {
		br.dead = true
	}

	if br.time > 255 {
		br.time = 0
	}

	br.mu.Unlock()
}

func (br *bird) paint() error {
	br.mu.RLock()

	rect := &sdl.Rect{
		X: br.posX,
		// TODO: Shouldn't this be using 0 .. n instead of -n .. n
		Y: ((ConstScreenHeight - int32(br.posY)) - br.height) / 2,
		W: br.width,
		H: br.height,
	}

	ind := br.time % uint8(len(br.textures))
	if err := br.ren.Copy(br.textures[ind], nil, rect); err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}

	br.mu.RUnlock()
	return nil
}

func (br *bird) touch(pp *pipe) {
	br.mu.RLock()
	pp.mu.RLock()

	// TODO: Better presition on pipe collition
	fmt.Printf("BirdX %d - PipeX %d | BirdY %d - Pipe Y %d\n", br.posX, pp.posX, -br.posY, pp.posY)
	if br.posX >= pp.posX && br.posX <= pp.posX+pp.width {
		if -br.posY <= pp.posY {
			br.dead = true
		}
	}

	pp.mu.RUnlock()
	br.mu.RUnlock()
}

func (br *bird) jump() {
	br.mu.Lock()
	if isFirstJump {
		br.speed = -constSpeed
	}
	isFirstJump = !isFirstJump
	br.mu.Unlock()
}

func (br *bird) isDead() bool {
	return br.dead == true
}

func (br *bird) restart() {
	br.mu.Lock()
	br.posY = ConstBirdStartingPoint
	br.dead = false
	br.time = 0
	br.speed = 0
	br.mu.Unlock()
}

func (br *bird) destroy() {
	br.mu.Lock()
	for _, t := range br.textures {
		t.Destroy()
	}
	br.mu.Unlock()
}
