package character

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/grid"
)

type Player struct {
	currentPosition *rl.Vector2
	targetPosition  *rl.Vector2
	moving          bool
	speed           float32 // Using rl.GetFrameTime will result in a pixel per second paradigm
	spriteAsset     string
	texture         *rl.Texture2D
	mapId           int
}

func (p *Player) ReachDestination() bool {
	if p.currentPosition.X == p.targetPosition.X && p.currentPosition.Y == p.targetPosition.Y {
		return true
	}
	return false
}

func (p *Player) GetCurrentPosition() *rl.Vector2 {
	return p.currentPosition
}

func (p *Player) SetCurrentPosition(position *rl.Vector2) {
	p.currentPosition = position
}

func (p *Player) GetTargetPosition() *rl.Vector2 {
	return p.targetPosition
}

func (p *Player) SetTargetPosition(target *rl.Vector2) {
	p.moving = true
	p.targetPosition = target
}

func (p *Player) IsMoving() bool {
	return p.moving
}

// Stop moving and set the current position to the target position
// A player force-stop requires a call to SetTargetPosition() before StopMoving()
func (p *Player) StopMoving() {
	p.currentPosition = p.targetPosition
	p.moving = false
}

// NOTE
// should eventually be orchestrated via
// Update function, which can be part of Entity interface
// Function to update player position,
func (p *Player) MoveForFrame() {
	time := rl.GetFrameTime()
	p.currentPosition.X = grid.ApproachPoint(p.currentPosition.X, p.targetPosition.X, p.speed*time)
	p.currentPosition.Y = grid.ApproachPoint(p.currentPosition.Y, p.targetPosition.Y, p.speed*time)
}

// Get the rectangle with the character sprite
// For now this is basically the initial character based on metadata
// which we are hardcoding in here for now
func (p *Player) GetSprite() (*rl.Texture2D, *rl.Rectangle) {
	if p.texture == nil {
		texture := rl.LoadTexture(p.spriteAsset)
		p.texture = &texture
	}

	var size float32 = 96 // Hardcoded for now
	textureRectangle := rl.NewRectangle(48, 48, size, size)

	return p.texture, &textureRectangle
}

func (p *Player) Destroy() {
	p.unloadSprite()
}

func (p *Player) unloadSprite() {
	if p.texture != nil {
		rl.UnloadTexture(*p.texture)
	}
}

func (p *Player) GetCurrentMap() int {
	// TODO
	// Implement persistent player state
	return 1000
}

func NewPlayer(spawn *rl.Vector2, speed float32) Entity {
	return &Player{
		currentPosition: spawn,
		speed:           speed,
		moving:          false,
		spriteAsset:     "./assets/TinySwordsFreePack/Units/Blue Units/Warrior/Warrior_Idle.png",
	}
}
