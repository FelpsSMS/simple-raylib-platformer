package main

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sprite struct {
	FrameCounter int
	CurrentFrame int
	FrameSpeed   int
	FrameRec     rl.Rectangle
	Position     rl.Vector2
	Texture      rl.Texture2D
}

type State int

const (
	IDLE State = iota
	MOVING
	PAUSED
)

type Player struct {
	X         float32
	Y         float32
	Width     float32
	Height    float32
	Hitbox    rl.Rectangle
	State     State
	Sprite    Sprite
	RightSide bool
	isJumping bool
	isFalling bool
	originalY float32
}

type OffsetParams struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

var (
	instance *Player
	once     sync.Once
)

func GetPlayer() *Player {
	once.Do(func() {
		playerTexture := rl.LoadTexture("assets/player/player-spritemap-v9.png")

		instance = &Player{
			State:     IDLE,
			X:         100,
			Y:         100,
			Width:     16,
			Height:    38,
			RightSide: true,
			isJumping: false,
			originalY: 100,
			isFalling: true,
			Sprite: Sprite{
				Texture:  playerTexture,
				FrameRec: rl.NewRectangle(0, 0, float32(playerTexture.Width/8), float32(playerTexture.Height/4)),
				Position: rl.Vector2{X: 100, Y: 100},
			},
		}
	})
	return instance
}
func (p *Player) CheckForPause() {

	if rl.IsKeyPressed(rl.KeyP) {

		if p.State != PAUSED {
			p.State = PAUSED

			for _, mob := range mobs {
				mob.State = PAUSED
			}

			return
		}

		for _, mob := range mobs {
			mob.State = MOVING
		}

		p.State = IDLE
	}
}

func (p *Player) CheckForMovement() {
	speed := float32(4)
	p.State = IDLE

	if rl.IsKeyPressed(rl.KeyUp) && !p.isJumping && !p.isFalling {
		p.isJumping = true
		p.originalY = p.Y
	}

	for _, tile := range tiles {

		if !tile.solid {
			continue
		}

		offsetSide := float32(1)

		if !p.RightSide {
			offsetSide = -1
		}

		if (rl.CheckCollisionRecs(p.OffsetHitbox(OffsetParams{X: speed * 1.1 * offsetSide}), tile.Hitbox)) {
			p.X -= 0.5 * offsetSide
			return
		}
	}

	if rl.IsKeyDown(rl.KeyRight) {
		p.X += speed
		p.RightSide = true
		p.State = MOVING
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		p.X -= speed
		p.RightSide = false
		p.State = MOVING
	}
}

func (p *Player) ApplyGravity() {
	gravity := float32(-1)
	jumpHeight := float32(100)

	for _, tile := range tiles {

		if !tile.solid {
			continue
		}

		if p.originalY-p.Y < jumpHeight && p.isJumping {
			gravity = 2

			if rl.CheckCollisionRecs(p.OffsetHitbox(OffsetParams{Y: gravity * -1.1}), tile.Hitbox) {
				p.isFalling = true
				p.isJumping = false
				return
			}

		} else {
			p.isJumping = false
		}

		if rl.CheckCollisionRecs(p.OffsetHitbox(OffsetParams{Y: gravity * -1.1}), tile.Hitbox) {
			p.isFalling = false
			return
		}
	}

	p.Y -= gravity
	p.isFalling = true
}

func (p *Player) Draw() {
	rect := rl.Rectangle{X: p.X, Y: p.Y, Width: p.Width, Height: p.Height}
	p.Hitbox = rect
	rl.DrawRectangleRec(rect, rl.DarkBlue)

	if p.isFalling || p.isJumping {
		p.Sprite.FrameRec = rl.NewRectangle(float32(p.Sprite.Texture.Width/8)*7, 0, float32(p.Sprite.Texture.Width/8), float32(p.Sprite.Texture.Height/4))

	} else {
		if p.State == IDLE {
			p.Sprite.FrameRec = rl.NewRectangle(0, 0, float32(p.Sprite.Texture.Width/8), float32(p.Sprite.Texture.Height/4))
		}

		if p.State == MOVING {
			p.Sprite.FrameSpeed++

			if p.Sprite.FrameSpeed > 8 {
				p.Sprite.FrameRec = rl.NewRectangle(float32(p.Sprite.Texture.Width/8)*float32(p.Sprite.FrameCounter), float32(p.Sprite.Texture.Height/4)*3, float32(p.Sprite.Texture.Width/8), float32(p.Sprite.Texture.Height/4))
				p.Sprite.FrameCounter++

				if p.Sprite.FrameCounter >= 6 {
					p.Sprite.FrameCounter = 0
				}

				p.Sprite.FrameSpeed = 0
			}
		}
	}

	if p.RightSide {
		if p.Sprite.FrameRec.Width < 0 {
			p.Sprite.FrameRec.Width *= -1
		}
	} else {
		if p.Sprite.FrameRec.Width > 0 {
			p.Sprite.FrameRec.Width *= -1
		}
	}

	rl.DrawTextureRec(p.Sprite.Texture, p.Sprite.FrameRec, rl.Vector2{X: p.X - p.Hitbox.Width, Y: p.Y - p.Hitbox.Height/4}, rl.White)
}

func (p *Player) OffsetHitbox(offset OffsetParams) rl.Rectangle {
	return rl.Rectangle{X: p.Hitbox.X + offset.X, Y: p.Hitbox.Y + offset.Y, Width: p.Hitbox.Width + offset.Width, Height: p.Hitbox.Height + offset.Height}
}
