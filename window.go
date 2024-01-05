package main

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	box        rl.Rectangle
	sprite     Sprite
	components []*Component
}

type Component struct {
	box    rl.Rectangle
	sprite Sprite
}

var (
	window        *Window
	inventoryOnce sync.Once
)

func GetInventoryWindow() *Window {
	rect := rl.Rectangle{X: float32(SCREEN_WIDTH / 4), Y: float32(SCREEN_HEIGHT / 4), Width: 400, Height: 200}

	inventoryOnce.Do(func() {
		window = &Window{
			box: rect,
		}
		// window.box.X = float32(SCREEN_WIDTH / 4)
		// window.box.Y = float32(SCREEN_HEIGHT / 4)
		// window.box.Width = 400
		// window.box.Height = 200
	})

	return window
}
