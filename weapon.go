package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Weapon struct {
	Damage float32
	Hitbox rl.Rectangle
	Name   string
}
