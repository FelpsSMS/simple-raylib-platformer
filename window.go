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
	box     rl.Rectangle
	window  *Window
	sprite  Sprite
	text    string
	context interface{}
	onClick []func()
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

func (window *Window) SetWindowIsOpen(isOpen bool) {
	window.isOpen = isOpen

	index := FindElementIndex(openWindows, inventoryWindow)

	if window.isOpen {
		if index == -1 {
			openWindows = append(openWindows, inventoryWindow)

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
	currentWindow := component.window

windowLoop:
	for {

		for _, window := range openWindows {

			if window.parent == currentWindow {
				index := FindElementIndex(openWindows, currentWindow)

				if index != -1 {
					openWindows = RemoveFromSlice(openWindows, index)
				}

				currentWindow = window

				continue windowLoop
			}

			break windowLoop
		}
	}
}

func (component *Component) Draw() {
	//rl.DrawRectangleRec(component.box, rl.Orange)
}

/* func (component *Component) CheckComponentsIfWindowIsOpen() {
	if component.window.isOpen {

		for _, windowComponent := range openComponents {
			index := FindElementIndex(openComponents, windowComponent)

			if index == -1 {
				openComponents = append(openComponents, windowComponent)
			}
		}

	} else {

		for _, windowComponent := range openComponents {
			index := FindElementIndex(openComponents, windowComponent)

			if index != -1 {
				openComponents = RemoveFromSlice(openComponents, index)
			}
		}
	}
} */

func (component *Component) CheckForTogglingItemWindow() {
	//player := GetPlayer()
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseButtonRight) && rl.CheckCollisionPointRec(mousePos, component.box) {
		logger.Print("open description window")
		/* index := -1

		if item, ok := component.context.(*Item); ok {
			index = FindElementIndex(player.Inventory, item)
		}

		if index != -1 {
			player.Inventory = RemoveFromSlice(player.Inventory, index)
		} */
	}
}

func (component *Component) CheckForClickEvent() {
	for _, fn := range component.onClick {
		fn()
	}

}
