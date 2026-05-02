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
}

func (p *Player) HasArrived() bool {
	if p.currentPosition.X == p.targetPosition.X && p.currentPosition.Y == p.targetPosition.Y {
		return true
	}
	return false
}

func (p *Player) GetCurrentPosition() *rl.Vector2 {
	return p.currentPosition
}

func (p *Player) SetCurrentPosition(position *rl.Vector2) {
	p.currentPosition.X = position.X
	p.currentPosition.Y = position.Y
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

func (p *Player) StopMoving() {
	p.currentPosition = p.targetPosition
	p.moving = false
}

func (p *Player) UpdatePosition(time float32) {
	p.currentPosition.X = grid.ApproachPoint(p.currentPosition.X, p.targetPosition.X, p.speed*time)
	p.currentPosition.Y = grid.ApproachPoint(p.currentPosition.Y, p.targetPosition.Y, p.speed*time)
}

func (p *Player) LoadSprite() Entity {
	texture := rl.LoadTexture(p.spriteAsset)
	p.texture = &texture
	return p
}

func (p *Player) UnloadSprite() {
	if p.texture != nil {
		rl.UnloadTexture(*p.texture)
	}
}

func (p *Player) Draw() {
	if p.texture == nil {
		p.LoadSprite()
	}

	playerPosition := p.GetCurrentPosition()

	rl.DrawTexture(
		*p.texture,
		int32(playerPosition.X)-p.texture.Width/2,
		int32(playerPosition.Y)-(p.texture.Height-6),
		rl.White,
	)
}

func NewPlayer(spawn *rl.Vector2, speed float32) Entity {
	return &Player{
		currentPosition: spawn,
		speed:           speed,
		moving:          false,
		spriteAsset:     "assets/character/playerPlaceHolderSprite.png",
	}
}
