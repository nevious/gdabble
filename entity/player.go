package entity

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/animation"
	"main/grid"
	"main/types"
)

type FaceDirection int

const (
	North FaceDirection = iota
	East
	South
	West
)

type Character struct {
	currentPosition *rl.Vector2
	targetPosition  *rl.Vector2
	speed           float32 // Using rl.GetFrameTime will result in a pixel per second paradigm
	mapId           int
	entityType      EntityType
	actionState     EntityMovementState
	animations      map[EntityMovementState]animation.Animation
	facing          FaceDirection
}

/* Public -------------------------------------------------------------------------------------- */

func (p *Character) ReachDestination() bool {
	if p.currentPosition.X == p.targetPosition.X && p.currentPosition.Y == p.targetPosition.Y {
		return true
	}
	return false
}

func (p *Character) GetCurrentPosition() *rl.Vector2 {
	return p.currentPosition
}

func (p *Character) GetTargetPosition() *rl.Vector2 {
	return p.targetPosition
}

func (p *Character) SetTargetPosition(target *rl.Vector2) {
	p.actionState = StateWalk
	p.targetPosition = target
}

// Stop moving and set the current position to the target position
// A player force-stop requires a call to SetTargetPosition() before StopMoving()
func (p *Character) StopMoving() {
	p.currentPosition = p.targetPosition
	p.actionState = StateIdle
}

// Update the character state
func (p *Character) Update() {
	if p.GetActionState() == StateWalk {
		if p.ReachDestination() {
			p.StopMoving()
			return
		}

		t := rl.GetFrameTime()
		p.currentPosition.X = grid.ApproachPoint(p.currentPosition.X, p.targetPosition.X, p.speed*t)
		p.currentPosition.Y = grid.ApproachPoint(p.currentPosition.Y, p.targetPosition.Y, p.speed*t)

		if p.currentPosition.X > p.targetPosition.X {
			p.facing = West
		} else if p.currentPosition.X < p.targetPosition.X {
			p.facing = East
		}
	}
}

func (p *Character) GetSprite() *types.RenderItem {
	reverse := p.facing == West
	renderItem := p.animations[p.GetActionState()].NextFrame(reverse)
	return &renderItem
}

func (p *Character) Destroy() {}

func (p *Character) GetActionState() EntityMovementState {
	return p.actionState
}

func (p *Character) GetEntityType() EntityType {
	return p.entityType
}

func (p *Character) GetCurrentMap() int {
	// TODO Implement persistent player state
	return 1000
}

/* Private -------------------------------------------------------------------------------------- */

/* Constructor Method --------------------------------------------------------------------------- */

func NewPlayer(spawn *rl.Vector2, speed float32) Entity {
	// TODO
	// Construct this from persistant data source
	var anim = map[EntityMovementState]animation.Animation{
		StateIdle: animation.NewSpriteAnimation(
			"./assets/TinySwordsFreePack/Units/Blue Units/Warrior/Warrior_Idle.png",
			8, 0, 192, 8, 3, 2,
		),
		StateWalk: animation.NewSpriteAnimation(
			"./assets/TinySwordsFreePack/Units/Blue Units/Warrior/Warrior_Run.png",
			6, 0, 192, 8, 3, 2,
		),
	}

	return &Character{
		currentPosition: spawn,
		speed:           speed,
		actionState:     StateIdle,
		entityType:      PlayerType,
		animations:      anim,
		facing:          East,
	}
}

func NewEnemy(spawn rl.Vector2) Entity {
	// Enemies should also come from a persistent data source
	var anim = map[EntityMovementState]animation.Animation{
		StateIdle: animation.NewSpriteAnimation(
			"assets/TinySwordsFreePack/Units/Red Units/Pawn/Pawn_Idle Axe.png",
			8, 0, 192, 8, 3, 2,
		),
		StateWalk: animation.NewSpriteAnimation(
			"assets/TinySwordsFreePack/Units/Red Units/Pawn/Pawn_Run Axe.png",
			6, 0, 192, 8, 3, 2,
		),
		StateAttack: animation.NewSpriteAnimation(
			"assets/TinySwordsFreePack/Units/Red Units/Pawn/Pawn_Interact Axe.png",
			6, 0, 192, 8, 3, 2,
		),
	}

	return &Character{
		currentPosition: &spawn,
		speed:           50,
		entityType:      EnemeyType,
		actionState:     StateIdle,
		animations:      anim,
		facing:          West,
	}
}
