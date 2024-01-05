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

func (item *Item) Use() {
	itemMap[item.itemId]()
}

func heal(value float32) {
	player := GetPlayer()

	player.HP += value
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
