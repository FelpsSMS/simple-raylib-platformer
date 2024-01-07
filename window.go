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

func (window *Window) CheckForDrag(disableDrag bool) {
	if disableDrag {
		disableDragCounter = 1
	}

	if disableDragCounter < 60 && disableDragCounter > 0 {
		disableDragCounter++
		return

	} else {
		disableDragCounter = 0
	}

	bufferZone := float32(15.0)
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {

		if rl.CheckCollisionPointRec(mousePos, window.box) ||
			rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(window.box.X-bufferZone, window.box.Y-bufferZone, window.box.Width+bufferZone*2, window.box.Height+bufferZone*2)) {

			if !isDragging {
				isDragging = true
				dragOffset = rl.NewVector2(mousePos.X-window.box.X, mousePos.Y-window.box.Y)
			}

			window.box.X = mousePos.X - dragOffset.X
			window.box.Y = mousePos.Y - dragOffset.Y
		}
	} else {
		isDragging = false
	}
}

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
