package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MovePattern int

const (
	STILL MovePattern = iota
	FIXED_HORIZONTAL
	JUMPING
)

type AttackPattern int

const (
	PACIFIST AttackPattern = iota
	RANGED_BOTH_SIDES
	RANGED_FRONT
	MELEE_FRONT
	MELE_BOTH_SIDES
)

type Mob struct {
	X               float32
	Y               float32
	Width           float32
	Height          float32
	Hitbox          rl.Rectangle
	State           State
	Sprite          Sprite
	RightSide       bool
	isJumping       bool
	isFalling       bool
	originalY       float32
	HP              float32
	MaxHP           float32
	MovePattern     MovePattern
	SpawnX          float32
	SpawnY          float32
	Name            string
	MoveSpeed       float32
	HPBar           rl.Rectangle
	originalHPWidth float32
	Damage          float32
	dropTable       []ItemDrop
	projectile      Projectile
	attackPattern   AttackPattern
	shootCD         float32
	shootCDCounter  float32
}

func Spawn(mob Mob) *Mob {
	mob.SpawnX = mob.X
	mob.SpawnY = mob.Y
	mob.originalY = mob.Y
	mob.isFalling = true
	mob.isJumping = false
	mob.State = MOVING
	mob.RightSide = false
	mob.MaxHP = mob.HP
	mob.HPBar = rl.NewRectangle(mob.X, mob.Y+mob.Height, 20, 4)
	mob.originalHPWidth = 20
	mob.shootCDCounter = 0

	return &mob
}

func (mob *Mob) Move() {
	maxDistanceX := float32(50)

	switch mob.MovePattern {
	// case STILL:

	case FIXED_HORIZONTAL:

		if mob.X >= mob.SpawnX+maxDistanceX {
			mob.RightSide = true
		}

		if mob.X <= mob.SpawnX-maxDistanceX {
			mob.RightSide = false
		}

		if mob.RightSide {
			mob.X -= mob.MoveSpeed

		} else {
			mob.X += mob.MoveSpeed
		}
	}
}

func (mob *Mob) ApplyGravity() {
	gravity := float32(-1)
	jumpHeight := float32(100)

	for _, tile := range tiles {

		if !tile.solid {
			continue
		}

		if mob.originalY-mob.Y < jumpHeight && mob.isJumping {
			gravity = 2

			if rl.CheckCollisionRecs(mob.OffsetHitbox(OffsetParams{Y: gravity * -1.1}), tile.Hitbox) {
				mob.isFalling = true
				mob.isJumping = false
				return
			}

		} else {
			mob.isJumping = false
		}

		if rl.CheckCollisionRecs(mob.OffsetHitbox(OffsetParams{Y: gravity * -1.1}), tile.Hitbox) {
			mob.isFalling = false
			return
		}
	}

	mob.Y -= gravity
	mob.isFalling = true
}

func (mob *Mob) OffsetHitbox(offset OffsetParams) rl.Rectangle {
	return rl.Rectangle{X: mob.Hitbox.X + offset.X, Y: mob.Hitbox.Y + offset.Y, Width: mob.Hitbox.Width + offset.Width, Height: mob.Hitbox.Height + offset.Height}
}

func (mob *Mob) dropItems() {
	roll := rand.Intn(100)

	for _, itemDrop := range mob.dropTable {

		if roll <= itemDrop.chance {
			rect := rl.Rectangle{X: mob.X + mob.Width/2, Y: mob.Y + mob.Height/2, Width: itemDrop.item.hitbox.Width, Height: itemDrop.item.hitbox.Height}

			itemDrop.item.hitbox = rect
			itemsInMap = append(itemsInMap, &itemDrop.item)
		}
	}
}

func (mob *Mob) Attack() {

	switch mob.attackPattern {
	/* 	case PACIFIST: */

	case RANGED_BOTH_SIDES:
		hitbox := rl.Rectangle{
			X:      mob.Hitbox.X,
			Y:      mob.Hitbox.Y + mob.Height/2,
			Width:  mob.projectile.width,
			Height: mob.projectile.height,
		}

		projectile := SpawnProjectile(Projectile{
			hitbox:       hitbox,
			name:         mob.projectile.name,
			damage:       mob.projectile.damage,
			speed:        mob.projectile.speed,
			isFromPlayer: false,
		})

		projectile.rightSide = mob.RightSide

		if mob.shootCDCounter <= 0 {
			mob.shootCDCounter++
			logger.Print(mob.shootCDCounter)
			projectilesInMap = append(projectilesInMap, &projectile)

		} else if mob.shootCDCounter >= mob.shootCD {
			mob.shootCDCounter = 0

		} else {
			mob.shootCDCounter++
		}

	}

}

func (mob *Mob) Draw() {
	rect := rl.Rectangle{X: mob.X, Y: mob.Y, Width: mob.Width, Height: mob.Height}

	mob.Hitbox = rect

	if mob.HP > 0 {
		mob.HPBar.Width = mob.originalHPWidth * (mob.HP / mob.MaxHP)

		mob.HPBar.X = mob.X
		mob.HPBar.Y = mob.Y - mob.Height/2

		rl.DrawRectangleRec(mob.HPBar, rl.Red)
		rl.DrawRectangleRec(rect, rl.DarkGreen)
	} else {
		index := FindElementIndex(mobs, mob)

		if index != -1 {
			mobs = RemoveFromSlice(mobs, index)
		}

		mob.dropItems()
	}

	// if mob.isFalling || mob.isJumping {
	// 	mob.Sprite.FrameRec = rl.NewRectangle(float32(mob.Sprite.Texture.Width/8)*7, 0, float32(mob.Sprite.Texture.Width/8), float32(mob.Sprite.Texture.Height/4))

	// } else {
	// 	if mob.State == IDLE {
	// 		mob.Sprite.FrameRec = rl.NewRectangle(0, 0, float32(mob.Sprite.Texture.Width/8), float32(mob.Sprite.Texture.Height/4))
	// 	}

	// 	if mob.State == MOVING {
	// 		mob.Sprite.FrameSpeed++

	// 		if mob.Sprite.FrameSpeed > 8 {
	// 			mob.Sprite.FrameRec = rl.NewRectangle(float32(mob.Sprite.Texture.Width/8)*float32(mob.Sprite.FrameCounter), float32(mob.Sprite.Texture.Height/4)*3, float32(mob.Sprite.Texture.Width/8), float32(mob.Sprite.Texture.Height/4))
	// 			mob.Sprite.FrameCounter++

	// 			if mob.Sprite.FrameCounter >= 6 {
	// 				mob.Sprite.FrameCounter = 0
	// 			}

	// 			mob.Sprite.FrameSpeed = 0
	// 		}
	// 	}
	// }

	// if mob.RightSide {
	// 	if mob.Sprite.FrameRec.Width < 0 {
	// 		mob.Sprite.FrameRec.Width *= -1
	// 	}
	// } else {
	// 	if mob.Sprite.FrameRec.Width > 0 {
	// 		mob.Sprite.FrameRec.Width *= -1
	// 	}
	// }

	// rl.DrawTextureRec(mob.Sprite.Texture, mob.Sprite.FrameRec, rl.Vector2{X: mob.X - mob.Hitbox.Width, Y: mob.Y - mob.Hitbox.Height/4}, rl.White)
}
