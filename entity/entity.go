package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nevious/gdabble/types"
)

type Entity interface {
	// Check if destination is reached
	ReachDestination() bool
	// Get the current Position in Pixel space
	GetCurrentPosition() *rl.Vector2
	// Get the target position in Pixel space
	GetTargetPosition() *rl.Vector2
	// Set the target position in Pixel space
	SetTargetPosition(*rl.Vector2)
	// Check the entity's movement state
	GetActionState() EntityMovementState
	// Stop moving and set the current position to the target position
	StopMoving()
	// Get this entities sprite texture
	GetSprite() *types.RenderItem
	// Update the entity state
	Update()
	// Destroy Entitiy and unload the sprite
	Destroy()
	// Get the current map of the entity
	GetCurrentMap() int
	// Get the entities Type
	GetEntityType() EntityType
}

type EntityType int

const (
	PlayerType EntityType = iota
	EnemeyType
	NpcType
)

var entitytTypeNames = map[EntityType]string{
	PlayerType: "Player",
	EnemeyType: "Enemy",
	NpcType:    "NPC",
}

func (e EntityType) String() string {
	return entitytTypeNames[e]
}
