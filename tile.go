package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Tile struct {
	X      float32
	Y      float32
	Size   float32
	Hitbox rl.Rectangle
	solid  bool
}

func CreateTile(x float32, y float32, solid bool) Tile {
	var tile Tile
	tile.X = x
	tile.Y = y
	tile.Size = 32
	tile.solid = solid

	return tile
}
