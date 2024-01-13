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
var projectilesInMap []*Projectile
var invulnerabilityTimer = 0
var isDragging = false
var dragOffset rl.Vector2
var disableDragCounter = 0

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

	//Disable esc key for closing the game
	rl.SetExitKey(0)

	player := GetPlayer()
	rl.SetTargetFPS(60)

	defer rl.UnloadTexture(player.Sprite.Texture)
	basicPotion := Item{name: "Basic potion", itemId: "basicPotion", hitbox: rl.NewRectangle(0, 0, 12, 12)}

	basicMob := Spawn(Mob{Name: "Test", X: 400, Y: 350, Width: 30, Height: 40, HP: 100, MoveSpeed: 2, MovePattern: FIXED_HORIZONTAL, Damage: 5})
	basicMob.dropTable = append(basicMob.dropTable, ItemDrop{item: basicPotion, chance: 100})

	basicRangedMob := Spawn(Mob{
		Name: "Ranged", X: 500, Y: 300, Width: 30, Height: 80, HP: 100,
		MoveSpeed: 8, MovePattern: FIXED_HORIZONTAL, Damage: 5, attackPattern: RANGED_BOTH_SIDES_RANDOM,
		shootCD: 10,
	})

	basicRangedMob.projectile = Projectile{
		name:   "Basic Mob Shot",
		damage: 15,
		speed:  3,
		width:  12,
		height: 18,
	}

	mobs = append(mobs, basicMob)
	mobs = append(mobs, basicRangedMob)

	isInventoryOpen := false

	startDebugPlayer()

	for !rl.WindowShouldClose() {
		player.CheckForPause()
		//logger.Print(projectilesInMap)

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
				mob.Attack()
			}

			for _, item := range itemsInMap {
				item.DetectCollisionWithItem()
			}

			for _, projectile := range projectilesInMap {
				projectile.Move()
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

		for _, projectile := range projectilesInMap {
			projectile.Draw()
		}

		if player.State != PAUSED && player.State != DEAD {
			if rl.IsKeyPressed(rl.KeyI) {

				if isInventoryOpen {
					isInventoryOpen = false

				} else {
					isInventoryOpen = true
				}
			}

			if rl.IsKeyPressed(rl.KeyEscape) && isInventoryOpen {
				isInventoryOpen = false
			}
		}

		if isInventoryOpen {
			drawInventoryWindow()
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
	// addTileToMap(18, 3, false, rl.Blue)

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

func drawInventoryWindow() {
	player := GetPlayer()
	inventoryWindow := GetInventoryWindow()
	disableDrag := false

	rl.DrawRectangleRec(inventoryWindow.box, rl.Beige)

	initialX := inventoryWindow.box.X + 10
	initialY := inventoryWindow.box.Y + 10

	for _, item := range player.Inventory {
		itemBox := rl.Rectangle{X: initialX, Y: initialY, Width: 20, Height: 20}
		item.windowBox = itemBox

		rl.DrawRectangleRec(itemBox, rl.Blue)
		initialX += 10
		initialY += 10
		rl.DrawText(item.name, int32(initialX), int32(initialY), 8, rl.Black)

		disableDrag = item.CheckForUse()
	}

	inventoryWindow.CheckForDrag(disableDrag)
}

func startDebugPlayer() {
	p := GetPlayer()

	basicWeapon := Weapon{
		Name:     "Basic Sword",
		Damage:   50,
		Hitbox:   rl.NewRectangle(p.X+p.Width, p.Y+p.Height/2, 20, 4),
		isRanged: false,
	}

	basicRangedWeapon := Weapon{
		Name:     "Basic Gun",
		Hitbox:   rl.NewRectangle(p.X+p.Width, p.Y+p.Height/2, 10, 4),
		isRanged: true,
	}

	basicProjectile := Projectile{
		name:   "Basic Bullet",
		damage: 25,
		speed:  2,
		width:  5,
		height: 5,
	}

	p.HPBar = rl.NewRectangle(p.X, p.Y+p.Height, 20, 4)
	p.originalHPWidth = 20

	p.Weapon = basicWeapon

	p.Weapon = basicRangedWeapon
	p.projectileQuantity = 100
	p.projectileSlot = basicProjectile
}
