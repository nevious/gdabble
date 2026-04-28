package main

/* TODO
 * - [X] Normalize coordinate space. We're a little in cell space and a little in pixel space
 * - [-] Implement movement queue, so user can click 2 or 3 times before movement is blocked [^1]
 * - [ ] Implement enemies as Entities
 * - [ ] Implement a spire animation
 *
 * [^1]: The queue prooved extremly finicky. Queueing and dequeing while maintaining data integrity turns out to be a bit
 *       of a nightmare, which is why i ditched it for a "redirect approach" by removing the block when moving.
 */

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/gridsystem"
	"main/utils"
	//rg "github.com/gen2brain/raylib-go/raygui"
)

/*
 * Basic definitions and configurations
 */
var (
	/* Colors */
	accentBlue     = rl.NewColor(114, 204, 245, 255) // #72ccEB
	mainBackground = rl.NewColor(36, 39, 46, 255)    // #24272e
	mainForeground = rl.NewColor(226, 226, 226, 255) // #e2e2e2

	/* Sizing */
	width  int32 = 800
	height int32 = 450

	/* Font */
	gameFont rl.Font

	/* Map */
	cellWidth int = 30
)

/*
 * Typedef + implementation
 */

// ----------------------------------------------------------------------------
type Screen interface {
	Update() Screen // Return active/selected Screen
	Draw()          // Refresh on tick
}

type Entity interface {
	hasArrived(cellWidth int) bool
	getCurrentPosition() *rl.Vector2
	setCurrentPosition(*rl.Vector2)
	getTargetPosition() *rl.Vector2
	setTargetPosition(*rl.Vector2)
	isMoving() bool
	stopMoving()
	updatePosition(timer int)
}

// ----------------------------------------------------------------------------
type Game struct {
	world   Screen
	choices []string
	active  int
}

func (g *Game) Update() Screen {
	if rl.IsKeyPressed(rl.KeyEnter) {
		return g.world
	}
	return g
}

func (s *Game) Draw() {
	// TODO: Finish menu screen
	t := "Main Menu"
	var size int32 = 25
	offset := utils.GetCenterForText(t, size, gameFont)
	rl.DrawText(t, offset, 100, size, accentBlue)
}

func NewGame(w *World) Screen {
	return &Game{
		choices: []string{"Start", "Credit", "Quit"},
		active:  0,
		world:   w,
	}
}

// ----------------------------------------------------------------------------

type Player struct {
	currentPosition *rl.Vector2
	targetPosition  *rl.Vector2
	moving          bool
	speed           float32 // Using rl.GetFrameTime will result in a pixel per second paradigm
}

func (p *Player) hasArrived(cellWidth int) bool {
	if p.currentPosition.X == p.targetPosition.X && p.currentPosition.Y == p.targetPosition.Y {
		return true
	}
	return false
}

func (p *Player) getCurrentPosition() *rl.Vector2 {
	return p.currentPosition
}

func (p *Player) setCurrentPosition(position *rl.Vector2) {
	p.currentPosition.X = position.X
	p.currentPosition.Y = position.Y
}

func (p *Player) getTargetPosition() *rl.Vector2 {
	return p.targetPosition
}

func (p *Player) setTargetPosition(target *rl.Vector2) {
	p.moving = true
	p.targetPosition = target
}

func (p *Player) isMoving() bool {
	return p.moving
}

func (p *Player) stopMoving() {
	p.currentPosition = p.targetPosition
	p.moving = false
}

func (p *Player) updatePosition(time float32) {
	p.currentPosition.X = gridsystem.ApproachPoint(p.currentPosition.X, p.targetPosition.X, p.speed*time)
	p.currentPosition.Y = gridsystem.ApproachPoint(p.currentPosition.Y, p.targetPosition.Y, p.speed*time)
}

// ----------------------------------------------------------------------------

type World struct {
	cellWidth int
	player    *Player
}

func (w *World) Update() Screen {
	// Block if character is mooving
	if w.player.isMoving() {
		// if player has arrived, stop the move
		if w.player.hasArrived(w.cellWidth) {
			w.player.stopMoving()
			return w
		}

		w.player.updatePosition(rl.GetFrameTime())
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		vec := rl.GetMousePosition()
		clickedCell := gridsystem.GetCellFromPixelPosition(&vec, w.cellWidth)
		w.player.setTargetPosition(gridsystem.GetCenterCellCoordinates(clickedCell, w.cellWidth))
		utils.LogPlayerTransition(rl.LogDebug, *w.player.getCurrentPosition(), *w.player.getTargetPosition())
	}

	return w
}

func (w *World) Draw() {
	// Some debugging data
	detail := fmt.Sprintf("%dx%d@%d FPS", rl.GetScreenWidth(), rl.GetScreenHeight(), rl.GetFPS())
	rl.DrawText(detail, 10, int32(rl.GetScreenHeight())-30, gameFont.BaseSize+10, mainForeground)

	rl.DrawCircle(0, 0, 15, rl.Red)
	rl.DrawCircle(int32(rl.GetScreenWidth()), 0, 15, rl.Blue)
	rl.DrawCircle(0, int32(rl.GetScreenHeight()), 15, rl.Purple)
	rl.DrawCircle(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), 15, rl.Orange)

	vec_coord := *gridsystem.GetCoordinateSystem(w.cellWidth)
	t := fmt.Sprintf("%0.f x %0.f", vec_coord.X, vec_coord.Y)
	rl.DrawText(t, utils.GetCenterForText(t, 20, gameFont), 150, 20, mainForeground)

	for x := 0; x < int(vec_coord.X); x++ {
		from_x, to_x := int32(x*w.cellWidth-1), int32(x*w.cellWidth)
		rl.DrawLine(from_x, 0, to_x, int32(rl.GetScreenHeight()), rl.Pink)
	}

	for y := 0; y < int(vec_coord.Y); y++ {
		from_y, to_y := int32(y*w.cellWidth-1), int32(y*w.cellWidth)
		rl.DrawLine(0, from_y, int32(rl.GetScreenWidth()), to_y, rl.Pink)
	}

	// Actual Game logic
	if w.player.getCurrentPosition() != nil {
		p_pos := w.player.getCurrentPosition()
		cx, cy := int32(p_pos.X), int32(p_pos.Y)
		rl.DrawCircle(cx, cy, 15, rl.White)
	}
}

// ----------------------------------------------------------------------------

/*
 * Entrypoint with setup
 */
func warmup() {
	rl.InitWindow(width, height, "A dabblin' game!")
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetTraceLogLevel(rl.LogAll)
	rl.SetTargetFPS(120)
	gameFont = rl.GetFontDefault()
}

func main() {
	warmup()
	defer rl.CloseWindow()
	rl.TraceLog(rl.LogDebug, fmt.Sprintf("Font: %+v", gameFont))

	spawn_x := float32(rl.GetScreenWidth() / 2)
	spawn_y := float32(rl.GetScreenHeight() / 2)
	// Going through GetCellFromPixelPosition and GetCenterCellCoordinates
	// ensures we're in the middle of the cell, which is not neccessairly the middle
	// of the screen.
	spawn := gridsystem.GetCenterCellCoordinates(
		gridsystem.GetCellFromPixelPosition(&rl.Vector2{X: spawn_x, Y: spawn_y}, cellWidth),
		cellWidth,
	)

	player := &Player{
		moving:          false,
		currentPosition: spawn,
		speed:           175,
	}

	game := NewGame(
		&World{
			player:    player,
			cellWidth: cellWidth,
		},
	)

	var currentScreen Screen = game

	for !rl.WindowShouldClose() {
		currentScreen = currentScreen.Update()

		rl.BeginDrawing()
		rl.ClearBackground(mainBackground)
		currentScreen.Draw()
		rl.EndDrawing()
	}

}
