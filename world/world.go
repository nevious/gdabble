package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/character"
)

type World struct {
	cellWidth int
	player    character.Entity
}

func (w *World) Update() {
	if w.player.IsMoving() {
		// if player has arrived, stop the move
		if w.player.HasArrived() {
			w.player.StopMoving()
		}

		w.player.UpdatePosition(rl.GetFrameTime())
	}
}

func (w *World) GetPlayer() character.Entity {
	return w.player
}

func (w *World) GetCellWidth() int {
	return w.cellWidth
}

func NewWorld(cellWidth int, player character.Entity) *World {
	return &World{
		player:    player,
		cellWidth: cellWidth,
	}
}
