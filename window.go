package main

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	id                 string
	box                rl.Rectangle
	sprite             Sprite
	parent             *Window
	components         []*Component
	isOpen             bool
	isDragging         bool
	disableDrag        bool
	disableDragCounter int32
	zIndex             int32
}

type Component struct {
	box          rl.Rectangle
	window       *Window
	sprite       *Sprite
	text         string
	windowOffset rl.Vector2
	newWindow    *Window
	context      interface{}
	onClick      []func(rl.Vector2)
}

var (
	window        *Window
	inventoryOnce sync.Once
)

func (window *Window) CheckForDrag() {
	if window.disableDrag {
		window.disableDragCounter = 1
	}

	if window.disableDragCounter < 60 && window.disableDragCounter > 0 {
		window.disableDragCounter++
		return

	} else {
		window.disableDragCounter = 0
		window.disableDrag = false
	}

	bufferZone := float32(15.0)
	mousePos := GetMousePosition()

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {

		if rl.CheckCollisionPointRec(mousePos, window.box) ||
			rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(window.box.X-bufferZone, window.box.Y-bufferZone, window.box.Width+bufferZone*2, window.box.Height+bufferZone*2)) {

			for _, openWindow := range openWindows {

				if openWindow.isDragging && openWindow != window {
					return
				}
			}

			maxZIndexWindow := window
			maxZIndex := window.zIndex

			for _, openWindow := range openWindows {
				if openWindow != window && (rl.CheckCollisionPointRec(mousePos, openWindow.box) ||
					rl.CheckCollisionPointRec(mousePos, rl.NewRectangle(openWindow.box.X-bufferZone, openWindow.box.Y-bufferZone, openWindow.box.Width+bufferZone*2, openWindow.box.Height+bufferZone*2))) {

					if openWindow.zIndex > maxZIndex {
						maxZIndex = openWindow.zIndex
						maxZIndexWindow = openWindow
					}
				}
			}

			if !maxZIndexWindow.isDragging {
				maxZIndexWindow.isDragging = true
				dragOffset = rl.NewVector2(mousePos.X-maxZIndexWindow.box.X, mousePos.Y-maxZIndexWindow.box.Y)
			}

			maxZIndexWindow.box.X = mousePos.X - dragOffset.X
			maxZIndexWindow.box.Y = mousePos.Y - dragOffset.Y
		}
	} else {
		for _, openWindow := range openWindows {
			openWindow.isDragging = false
		}
	}
}

func GetInventoryWindow() *Window {
	inventoryOnce.Do(func() {
		window = &Window{
			id:     "inventoryWindow",
			parent: nil,
			zIndex: 0,
		}

		window.box.Width = 400
		window.box.Height = 200
	})

	return window
}

func (window *Window) Draw() {

	for _, openWindow := range openWindows {

		if window.parent == openWindow && window.isOpen {
			rl.DrawRectangleRec(window.box, rl.DarkGray)
		}
	}
}

func (window *Window) SetWindowIsOpen(isOpen bool) {
	window.isOpen = isOpen

	index := FindElementIndex(openWindows, window)

	if window.isOpen {
		if index == -1 {
			openWindows = append(openWindows, window)
		} else {
			openWindows = RemoveFromSlice(openWindows, index)
		}
	} else {
		if index != -1 {
			openWindows = RemoveFromSlice(openWindows, index)
		}
	}

}

func (component *Component) CloseWindow(mousePos rl.Vector2) {

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mousePos, component.box) {
		currentWindow := component.window

	windowLoop:
		for {

			for _, window := range openWindows {

				if window.parent == currentWindow {
					window.SetWindowIsOpen(false)

					currentWindow = window

					continue windowLoop
				}

			}
			currentWindow.SetWindowIsOpen(false)
			component.window.SetWindowIsOpen(false)

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

func (component *Component) CheckForTogglingItemWindow(mousePos rl.Vector2) {

	if rl.IsMouseButtonPressed(rl.MouseButtonRight) && rl.CheckCollisionPointRec(mousePos, component.box) {
		if component.newWindow.isOpen {
			component.newWindow.SetWindowIsOpen(false)
			return
		}

		component.newWindow.box = rl.NewRectangle(component.box.X, component.box.Y-80, 100, 60)
		component.newWindow.SetWindowIsOpen(true)

	}
}

func (component *Component) CheckForClickEvent() {
	mousePos := GetMousePosition()

	for _, openWindow := range openWindows {

		if openWindow != component.window && rl.CheckCollisionPointRec(mousePos, openWindow.box) && openWindow.zIndex > component.window.zIndex {
			return
		}
	}

	for _, fn := range component.onClick {
		fn(mousePos)
	}

}
