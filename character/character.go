package character

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Entity interface {
	// Check if destination is reached
	ReachDestination() bool
	// Get the current Position in Pixel space
	GetCurrentPosition() *rl.Vector2
	// Set the current Position in Pixel space
	SetCurrentPosition(*rl.Vector2)
	// Get the target position in Pixel space
	GetTargetPosition() *rl.Vector2
	// Set the target position in Pixel space
	SetTargetPosition(*rl.Vector2)
	// Check if the entitiy is currently moving
	IsMoving() bool
	// Stop moving and set the current position to the target position
	StopMoving()
	// Move a frames width towards the target position
	MoveForFrame()
	// Get this entities sprite texture
	GetSprite() (*rl.Texture2D, *rl.Rectangle)
	// Destroy Entitiy and unload the sprite
	Destroy()
	// Get the current map of the entity
	GetCurrentMap() int
}
