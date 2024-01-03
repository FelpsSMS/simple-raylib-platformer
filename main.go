package main

import (
	"image/color"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SCREEN_WIDTH = 800
var SCREEN_HEIGHT = 450
var GAME_TITLE = "PLATFORMER"
var logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
var tiles []Tile
var mobs []*Mob
var invulnerabilityTimer = 0

func main() {
	rl.InitWindow(int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), GAME_TITLE)
	defer rl.CloseWindow()

	player := GetPlayer()
	rl.SetTargetFPS(60)

	defer rl.UnloadTexture(player.Sprite.Texture)

	mobs = append(mobs, Spawn(Mob{Name: "Test", X: 400, Y: 350, Width: 30, Height: 40, HP: 100, MoveSpeed: 2, MovePattern: FIXED_HORIZONTAL, Damage: 50}))

	for !rl.WindowShouldClose() {
		player.CheckForPause()

		if player.State != PAUSED && player.State != DEAD {
			player.CheckForMovement()
			player.CheckForAttack()
			player.ApplyGravity()
			player.CheckForCollision()

			if invulnerabilityTimer > 0 {
				invulnerabilityTimer--
			} else {
				player.isInvunerable = false
			}

			for _, mob := range mobs {
				mob.Move()
				mob.ApplyGravity()
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

	// logger.Printf("TILES IN X %v \n", tilesInX)
	// logger.Printf("TILES IN Y %v \n", tilesInY)

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
