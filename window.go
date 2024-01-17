package main

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	box        rl.Rectangle
	sprite     Sprite
	components []*Component
	parent     *Window
	isOpen     bool
}

type Component struct {
	box           rl.Rectangle
	window        *Window
	sprite        *Sprite
	text          string
	windowOffset  rl.Vector2
	newWindowOpen bool
	context       interface{}
	onClick       []func()
}

var (
	window        *Window
	inventoryOnce sync.Once
)

func FindWindowIndex(slice []*Window, window *Window) int {
	for index, elementInSlice := range slice {
		if *elementInSlice == *window {
			return index
		}
	}

	return -1
}

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
			box:    rect,
			parent: nil,
		}
		// window.box.X = float32(SCREEN_WIDTH / 4)
		// window.box.Y = float32(SCREEN_HEIGHT / 4)
		// window.box.Width = 400
		// window.box.Height = 200
	})

	return window
}

func (window *Window) Draw() {

	for _, openWindow := range openWindows {
		logger.Print(openWindows)

		if window.parent == openWindow && window.isOpen {
			rl.DrawRectangleRec(window.box, rl.DarkGray)
		}
	}
}

func (window *Window) SetWindowIsOpen(isOpen bool) {
	window.isOpen = isOpen

	index := FindWindowIndex(openWindows, window)

	if window.isOpen {
		if index == -1 {
			openWindows = append(openWindows, window)

			for _, component := range window.components {
				componentIndex := FindElementIndex(openComponents, component)

				if componentIndex == -1 {
					openComponents = append(openComponents, component)
				}
			}
		}
	} else {
		if index != -1 {
			openWindows = RemoveFromSlice(openWindows, index)

			for _, component := range window.components {
				componentIndex := FindElementIndex(openComponents, component)

				if componentIndex != -1 {
					openComponents = RemoveFromSlice(openComponents, componentIndex)
				}
			}
		}
	}

}

func (component *Component) CloseWindow() {
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mousePos, component.box) {
		currentWindow := component.window

	windowLoop:
		for {

			for _, window := range openWindows {
				//logger.Print(window.parent)

				if window.parent == currentWindow {
					window.SetWindowIsOpen(false)

					currentWindow = window

					continue windowLoop
				}

			}
			currentWindow.SetWindowIsOpen(false)

			break windowLoop
		}
	}
}

func (component *Component) Draw() {
	if component.sprite != nil {
		component.box.X = component.window.box.X + component.windowOffset.X
		component.box.Y = component.window.box.Y + component.windowOffset.Y

		rl.DrawRectangleRec(rl.NewRectangle(component.box.X, component.box.Y, component.box.Width, component.box.Height), rl.Orange)
	}
}

func (component *Component) CheckForTogglingItemWindow() {
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseButtonRight) && rl.CheckCollisionPointRec(mousePos, component.box) && !component.newWindowOpen {
		box := rl.NewRectangle(component.box.X, component.box.Y, 100, 60)
		newWindow := &Window{box: box, parent: inventoryWindow, sprite: Sprite{}, isOpen: false}
		newWindow.SetWindowIsOpen(true)
		component.newWindowOpen = true
	}
}

func (component *Component) CheckForClickEvent() {
	for _, fn := range component.onClick {
		fn()
	}

}
