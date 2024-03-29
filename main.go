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
var dragOffset rl.Vector2
var playerInstance *Player
var openWindows []*Window
var inventoryWindow *Window
var camera rl.Camera2D
var targetY float32
var targetX float32
var previousTargetX float32
var previousTargetY float32

const cameraThreshold = 100

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

func DrawOutlinedText(text string, posX int32, posY int32, fontSize int32, color rl.Color, outlineSize int32, outlineColor rl.Color) {
	rl.DrawText(text, posX-outlineSize, posY-outlineSize, fontSize, outlineColor)
	rl.DrawText(text, posX+outlineSize, posY-outlineSize, fontSize, outlineColor)
	rl.DrawText(text, posX-outlineSize, posY+outlineSize, fontSize, outlineColor)
	rl.DrawText(text, posX+outlineSize, posY+outlineSize, fontSize, outlineColor)
	rl.DrawText(text, posX, posY, fontSize, color)
}

func GetMousePosition() rl.Vector2 {
	mousePos := rl.GetMousePosition()
	worldMousePos := rl.GetScreenToWorld2D(mousePos, camera)

	rl.DrawCircle(int32(worldMousePos.X), int32(worldMousePos.Y), 2, rl.Violet)

	return worldMousePos
}

func main() {
	rl.InitWindow(int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), GAME_TITLE)
	defer rl.CloseWindow()

	//Disable esc key for closing the game
	rl.SetExitKey(0)

	playerInstance = startDebugPlayer()
	inventoryWindow = GetInventoryWindow()

	camera = rl.NewCamera2D(rl.NewVector2(float32(SCREEN_WIDTH)/2, float32(SCREEN_HEIGHT)/2), rl.NewVector2(playerInstance.X, playerInstance.Y), 0, 1)

	closeInventoryWindowComponent := &Component{
		window:       inventoryWindow,
		sprite:       &Sprite{},
		windowOffset: rl.Vector2{X: inventoryWindow.box.Width - 25, Y: 5},
	}

	closeInventoryWindowComponent.box = rl.NewRectangle(inventoryWindow.box.X, window.box.Y, 20, 20)

	closeInventoryWindowComponent.onClick = append(closeInventoryWindowComponent.onClick, closeInventoryWindowComponent.CloseWindow)

	inventoryWindow.components = append(inventoryWindow.components, closeInventoryWindowComponent)

	rl.SetTargetFPS(60)

	defer rl.UnloadTexture(playerInstance.Sprite.Texture)
	startDebugItemsAndMobs()

	inventoryWindow.isOpen = false

	for !rl.WindowShouldClose() {
		playerInstance.CheckForPause()
		//logger.Print(projectilesInMap)

		if rl.IsKeyPressed(rl.KeyR) {
			resetWorld()
		}

		if playerInstance.State != PAUSED && playerInstance.State != DEAD {
			playerInstance.CheckForMovement()
			playerInstance.CheckForAttack()
			playerInstance.ApplyGravity()
			playerInstance.CheckForCollision()

			if invulnerabilityTimer > 0 {
				invulnerabilityTimer--
			} else {
				playerInstance.isInvunerable = false
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

		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode2D(camera)

		targetY = float32(SCREEN_HEIGHT) / 2
		targetX = float32(SCREEN_WIDTH) / 2

		if previousTargetX == 0 {
			previousTargetX = targetX
		}

		if previousTargetY == 0 {
			previousTargetY = targetY
		}

		if playerInstance.Y < cameraThreshold {
			targetY = float32(SCREEN_HEIGHT) + (playerInstance.Y - float32(SCREEN_HEIGHT))
		}

		if playerInstance.Y > float32(SCREEN_HEIGHT) {
			targetY = float32(SCREEN_HEIGHT) + (playerInstance.Y - float32(SCREEN_HEIGHT))
		}

		if playerInstance.X > float32(SCREEN_WIDTH) {
			targetX = float32(SCREEN_WIDTH) + (playerInstance.X - float32(SCREEN_WIDTH))
		}

		if playerInstance.X < 0 {
			targetX = float32(SCREEN_WIDTH) + (playerInstance.X - float32(SCREEN_WIDTH))
		}

		resultY := rl.Vector2Lerp(rl.NewVector2(0, camera.Target.Y), rl.NewVector2(0, targetY), 0.05)
		resultX := rl.Vector2Lerp(rl.NewVector2(camera.Target.X, 0), rl.NewVector2(targetX, 0), 0.05)

		camera.Target = rl.NewVector2(resultX.X, resultY.Y)

		initiateLevel()
		playerInstance.Draw()

		for _, mob := range mobs {
			mob.Draw()
		}

		for _, item := range itemsInMap {
			item.Draw()
		}

		for _, projectile := range projectilesInMap {
			projectile.Draw()
		}

		index := FindElementIndex(openWindows, inventoryWindow)

		if playerInstance.State != PAUSED && playerInstance.State != DEAD {

			if rl.IsKeyPressed(rl.KeyI) {

				if inventoryWindow.isOpen {
					inventoryWindow.SetWindowIsOpen(false)

				} else {
					inventoryWindow.SetWindowIsOpen(true)
				}
			}

			if rl.IsKeyPressed(rl.KeyEscape) && inventoryWindow.isOpen {
				inventoryWindow.SetWindowIsOpen(false)
			}

		} else {
			displayText := ""
			openWindows = nil

			if playerInstance.State == DEAD {
				displayText = "GAME OVER"
			}

			if playerInstance.State == PAUSED {
				displayText = "PAUSED"
			}

			DrawOutlinedText(displayText, int32(SCREEN_WIDTH)/2, int32(SCREEN_HEIGHT)/2, 40, rl.LightGray, 2, rl.Black)

			if playerInstance.State == DEAD {
				DrawOutlinedText("PRESS 'R' TO TRY AGAIN", int32(SCREEN_WIDTH)/2-150, int32(SCREEN_HEIGHT)/2+100, 30, rl.LightGray, 2, rl.Black)
			}
		}

		if inventoryWindow.isOpen && index != -1 {
			drawInventoryWindow()
		}

		const lerpFactor = 0.05

		for _, openWindow := range openWindows {
			// Update window position based on camera movement

			lerpVector := rl.Vector2Lerp(rl.NewVector2(openWindow.box.X, openWindow.box.Y), rl.NewVector2(targetX-(openWindow.box.X-previousTargetX), targetY-(openWindow.box.Y-previousTargetY)), lerpFactor)
			openWindow.box.X = lerpVector.X
			openWindow.box.Y = lerpVector.Y

			previousTargetX = targetX
			previousTargetY = targetY

			openWindow.Draw()
			openWindow.CheckForDrag()

			for _, component := range openWindow.components {
				component.Draw()
				component.CheckForClickEvent()
			}
		}

		rl.EndDrawing()
	}
}

func resetWorld() {
	itemsInMap = nil
	projectilesInMap = nil
	tiles = nil
	mobs = nil

	playerInstance = startDebugPlayer()
	startDebugItemsAndMobs()
}

func startDebugItemsAndMobs() {
	//Basic potion
	basicPotion := Item{name: "Basic potion", itemId: "basicPotion", hitbox: rl.NewRectangle(0, 0, 12, 12)}

	basicPotion.itemComponent = &Component{window: inventoryWindow, context: basicPotion}
	basicPotion.itemComponent.onClick = append(basicPotion.itemComponent.onClick, basicPotion.itemComponent.CheckForTogglingItemWindow)

	newWindow := &Window{}
	newWindow.parent = inventoryWindow
	newWindow.sprite = Sprite{}
	newWindow.isOpen = false
	newWindow.id = "itemDescription"
	newWindow.box = rl.NewRectangle(basicPotion.itemComponent.box.X, basicPotion.itemComponent.box.Y-80, 100, 60)
	newWindow.zIndex = 1

	closeNewWindowComponent := &Component{
		window:       newWindow,
		sprite:       &Sprite{},
		windowOffset: rl.Vector2{X: newWindow.box.Width - 25, Y: 5},
	}

	closeNewWindowComponent.box = rl.NewRectangle(newWindow.box.X, window.box.Y-window.box.Width, 20, 20)

	closeNewWindowComponent.onClick = append(closeNewWindowComponent.onClick, closeNewWindowComponent.CloseWindow)

	newWindow.components = append(newWindow.components, closeNewWindowComponent)

	basicPotion.itemComponent.newWindow = newWindow

	//Basic mob
	basicMob := Spawn(Mob{Name: "Test", X: 400, Y: 350, Width: 30, Height: 40, HP: 100, MoveSpeed: 2, MovePattern: FIXED_HORIZONTAL, Damage: 5})
	basicMob.dropTable = append(basicMob.dropTable, ItemDrop{item: basicPotion, chance: 100})

	//Basic ranged mob
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

	GetPlayer().Inventory = append(GetPlayer().Inventory, &basicPotion)
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

	rl.DrawRectangleRec(inventoryWindow.box, rl.Beige)

	initialX := inventoryWindow.box.X + 10
	initialY := inventoryWindow.box.Y + 10

	for _, item := range player.Inventory {
		itemBox := rl.Rectangle{X: initialX, Y: initialY, Width: 20, Height: 20}
		item.itemComponent.box = itemBox

		itemComponentIndex := FindElementIndex(inventoryWindow.components, item.itemComponent)

		if itemComponentIndex == -1 {
			inventoryWindow.components = append(inventoryWindow.components, item.itemComponent)
		}

		rl.DrawRectangleRec(itemBox, rl.Blue)
		initialX += 10
		initialY += 10
		rl.DrawText(item.name, int32(initialX), int32(initialY), 8, rl.Black)

		item.CheckForUse()
	}
}

func startDebugPlayer() *Player {
	p := NewPlayer()

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

	return p
}
