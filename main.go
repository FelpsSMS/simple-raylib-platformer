package main

import (
	"image/color"
	"log"
	"os"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SCREEN_WIDTH = 800
var SCREEN_HEIGHT = 450
var GAME_TITLE = "PLATFORMER"
var logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
var tiles []Tile
var mobs []*Mob
var itemsInMap []*Item
var invulnerabilityTimer = 0

func FindElementIndex[T any](slice []T, element T) int {
	for index, elementInSlice := range slice {
		if reflect.DeepEqual(elementInSlice, element) {
			return index
		}
	}

	return -1
}

func RemoveFromSlice[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func main() {
	rl.InitWindow(int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), GAME_TITLE)
	defer rl.CloseWindow()

	player := GetPlayer()
	rl.SetTargetFPS(60)

	defer rl.UnloadTexture(player.Sprite.Texture)
	basicPotion := Item{name: "Basic potion", itemId: "basicPotion", hitbox: rl.NewRectangle(0, 0, 12, 12)}

	basicMob := Spawn(Mob{Name: "Test", X: 400, Y: 350, Width: 30, Height: 40, HP: 100, MoveSpeed: 2, MovePattern: FIXED_HORIZONTAL, Damage: 5})
	basicMob.dropTable = append(basicMob.dropTable, ItemDrop{item: basicPotion, chance: 60})

	mobs = append(mobs, basicMob)

	for !rl.WindowShouldClose() {
		player.CheckForPause()

		if player.State != PAUSED && player.State != DEAD {
			player.CheckForMovement()
			player.CheckForAttack()
			player.ApplyGravity()
			player.CheckForCollision()

			logger.Println(player.Inventory)

			if invulnerabilityTimer > 0 {
				invulnerabilityTimer--
			} else {
				player.isInvunerable = false
			}

			for _, mob := range mobs {
				mob.Move()
				mob.ApplyGravity()
			}

			for _, item := range itemsInMap {
				item.DetectCollisionWithItem()
			}

		} else {
			displayText := ""

			if player.State == DEAD {
				displayText = "GAME OVER"
			}

			if player.State == PAUSED {
				displayText = "PAUSED"
			}

			rl.DrawText(displayText, int32(SCREEN_WIDTH)/2, int32(SCREEN_HEIGHT)/2, 40, rl.LightGray)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		initiateLevel()
		player.Draw()

		for _, mob := range mobs {
			mob.Draw()
		}

		for _, item := range itemsInMap {
			item.Draw()
		}

		rl.EndDrawing()
	}
}

func initiateLevel() {
	tileSize := 32
	tileXCoord := 0
	tilesInX := SCREEN_WIDTH / tileSize
	tilesInY := SCREEN_HEIGHT / tileSize

	tiles = nil

	for tileXCoord < tilesInX {
		addTileToMap(tileXCoord, 1, true, rl.Brown)

		tileXCoord += 1
	}

	addTileToMap(4, 2, true, rl.Brown)
	addTileToMap(5, 2, true, rl.Brown)
	addTileToMap(6, 2, true, rl.Brown)
	addTileToMap(7, 2, true, rl.Brown)

	addTileToMap(4, 8, true, rl.Brown)
	addTileToMap(8, 5, true, rl.Brown)
	addTileToMap(18, 3, false, rl.Blue)

	for i := 2; i <= tilesInY; i++ {
		addTileToMap(tilesInX-1, i, true, rl.DarkPurple)
	}

}

func addTileToMap(x int, y int, solid bool, color color.RGBA) {
	tile := CreateTile(float32(x), float32(y), solid)

	rect := rl.Rectangle{X: float32(tile.X) * tile.Size, Y: float32(SCREEN_HEIGHT) - tile.Y*tile.Size, Width: float32(tile.Size), Height: float32(tile.Size)}
	tile.Hitbox = rect

	rl.DrawRectangleRec(rect, color)

	tiles = append(tiles, tile)
}
