package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ItemType int

var itemMap = map[string]func(){
	"basicPotion": func() { heal(20) },
}

const (
	CONSUMABLE ItemType = iota
	EQUIPMENT
	MISCELLANEOUS
)

type ItemDrop struct {
	item   Item
	chance int
}

type Item struct {
	itemId    string
	itemType  ItemType
	name      string
	sprite    Sprite
	hitbox    rl.Rectangle
	windowBox rl.Rectangle
}

func (item *Item) CheckForUse() bool {
	player := GetPlayer()
	mousePos := rl.GetMousePosition()
	disableDrag := false

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(mousePos, item.windowBox) {
		disableDragCounter = 1
		itemMap[item.itemId]()

		index := FindElementIndex(player.Inventory, item)

		if index != -1 {
			player.Inventory = RemoveFromSlice(player.Inventory, index)
		}
	}

	return disableDrag
}

func heal(value float32) {
	player := GetPlayer()

	if player.HP+value <= player.MaxHP {
		player.HP += value
	} else {
		player.HP = player.MaxHP
	}
}

func (item *Item) DetectCollisionWithItem() {
	player := GetPlayer()

	if rl.CheckCollisionRecs(player.Hitbox, item.hitbox) {
		index := FindElementIndex(itemsInMap, item)

		if index != -1 {
			itemsInMap = RemoveFromSlice(itemsInMap, index)
		}

		player.Inventory = append(player.Inventory, item)
	}
}

func (item *Item) Draw() {
	rl.DrawRectangleRec(item.hitbox, rl.DarkPurple)
}
