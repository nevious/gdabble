package character

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Entity interface {
	HasArrived() bool
	GetCurrentPosition() *rl.Vector2
	SetCurrentPosition(*rl.Vector2)
	GetTargetPosition() *rl.Vector2
	SetTargetPosition(*rl.Vector2)
	IsMoving() bool
	StopMoving()
	UpdatePosition(timer float32)
	GetCharacterSprite() (*rl.Texture2D, *rl.Rectangle)
	DestroyCharacter()
	GetCurrentMap() int
}
