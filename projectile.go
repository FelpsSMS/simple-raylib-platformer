package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Projectile struct {
	name         string
	hitbox       rl.Rectangle
	sprite       Sprite
	damage       float32
	isFromPlayer bool
	speed        float32
	width        float32
	height       float32
	rightSide    bool
}

func SpawnProjectile(projectile Projectile) Projectile {
	var newProjectile Projectile
	newProjectile.hitbox = projectile.hitbox
	newProjectile.damage = projectile.damage
	newProjectile.speed = projectile.speed
	newProjectile.isFromPlayer = projectile.isFromPlayer

	return newProjectile
}

func (projectile *Projectile) Move() {
	if checkForProjectileCollisions(projectile) {
		return
	}

	for _, tile := range tiles {

		if rl.CheckCollisionRecs(projectile.hitbox, tile.Hitbox) {
			index := FindElementIndex(projectilesInMap, projectile)

			if index != -1 {
				projectilesInMap = RemoveFromSlice(projectilesInMap, index)
			}
		}
	}

	if projectile.rightSide {
		projectile.hitbox.X += projectile.speed

	} else {
		projectile.hitbox.X -= projectile.speed
	}
}

func checkForProjectileCollisions(projectile *Projectile) bool {
	player := GetPlayer()

	if projectile.isFromPlayer {

		for _, mob := range mobs {

			if rl.CheckCollisionRecs(mob.Hitbox, projectile.hitbox) {
				index := FindElementIndex(projectilesInMap, projectile)

				if index != -1 {
					projectilesInMap = RemoveFromSlice(projectilesInMap, index)
				}

				mob.HP -= projectile.damage
				return true
			}
		}
	} else if rl.CheckCollisionRecs(player.Hitbox, projectile.hitbox) {
		index := FindElementIndex(projectilesInMap, projectile)

		if index != -1 {
			projectilesInMap = RemoveFromSlice(projectilesInMap, index)
		}

		player.HP -= projectile.damage
		return true
	}

	return false
}

func (projectile *Projectile) Draw() {
	rl.DrawRectangleRec(projectile.hitbox, rl.Blue)
}
