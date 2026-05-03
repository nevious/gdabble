package main

/* TODO
 * - [X] Normalize coordinate space. We're a little in cell space and a little in pixel space
 * - [-] Implement movement queue, so user can click 2 or 3 times before movement is blocked [^1]
 * - [ ] Implement enemies as Entities
 * - [X] Decouple World and Menu Screens
 * - [ ] Implement a sprite animation
 * - [ ] Player spwn is in limbo if the map is not big enough
 *
 * [^1]: The queue prooved extremly finicky. Queueing and dequeing while maintaining data integrity turns out to be a bit
 *       of a nightmare, which is why i ditched it for a "redirect approach" by removing the block when moving.
 */

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/character"
	"main/grid"
	"main/ui"
	"main/utils"
	"main/world"
)

/*
 * Basic definitions and configurations
 * These could be moved into a config file.
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
	cellWidth float32 = 32

	/* Default speed */
	speed float32 = 175
)

// Create and spawn a player in the middle of the initial grid
func createPlayer(speed float32) character.Entity {
	spawnX := float32(rl.GetScreenWidth() / 2)
	spawnY := float32(rl.GetScreenHeight() / 2)

	// Going through GetCellFromPixelPosition and GetCenterCellCoordinates
	// ensures we're in the middle of the cell, which is not neccessairly the middle
	// of the screen.
	spawn := grid.GetCenterCellCoordinates(
		grid.GetCellFromPixelPosition(&rl.Vector2{X: spawnX, Y: spawnY}, cellWidth),
		cellWidth,
	)

	return character.NewPlayer(spawn, speed)
}

// Create a world with a player
func createWorld(cellWidth float32, player character.Entity) *world.World {
	return world.NewWorld(cellWidth, player)
}

// Create an empty Credits SCreen
func createCreditScreen(font *rl.Font, color, highlight rl.Color) ui.Screen {
	return ui.NewCreditScreen(font, color, highlight)
}

// Create a gameScren that contains the world and displays it
func createGameScreen(gameWorld *world.World, font *rl.Font, color, highlight rl.Color, camera *rl.Camera2D) ui.Screen {
	return ui.NewGameScreen(gameWorld, font, color, highlight, camera)
}

// Create a simply quick screen
func createQuitScreen() ui.Screen {
	return &ui.QuitScreen{}
}

// Create an empty menuscreen
func createMenuScreen(font *rl.Font, color, highlight rl.Color, items *[]ui.MenuItem) ui.Screen {
	return ui.NewMenuScreen(font, color, highlight, items)
}

// Create the settings object and define the settings displayed in it
func createSettings(font *rl.Font, color, hightlight rl.Color) ui.Screen {
	var values *[]ui.Setting
	return ui.NewSettingsScreen(values, font, color, hightlight)
}

// Create a camera and point it at the player
func createCamera(player character.Entity) *rl.Camera2D {
	return &rl.Camera2D{
		Target: *player.GetCurrentPosition(),
		Zoom:   1, // defaults to 0, infinite zoom
	}
}

func initRaylib() {
	//rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(width, height, "A dabblin' game!")
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetTraceLogLevel(rl.LogAll)
	rl.SetTargetFPS(120)
	gameFont = rl.GetFontDefault()
	utils.LogDebug("Font: %+v", gameFont)
}

func main() {
	initRaylib()
	defer rl.CloseWindow()

	// Build things in a constructor way
	// Might not be needed everywhere, but should give us options
	player := createPlayer(speed)
	defer player.DestroyCharacter()
	camera := createCamera(player)
	world := createWorld(cellWidth, player)
	gameScreen := createGameScreen(world, &gameFont, mainForeground, accentBlue, camera)
	creditScreen := createCreditScreen(&gameFont, mainForeground, accentBlue)
	settingsScreen := createSettings(&gameFont, mainForeground, accentBlue)
	quitScreen := createQuitScreen()

	// Create the menuItems
	menuItems := []ui.MenuItem{
		{Label: "Start Game", Screen: gameScreen},
		{Label: "Credit", Screen: creditScreen},
		{Label: "Settings", Screen: settingsScreen},
		{Label: "Quit", Screen: quitScreen},
	}
	menuScreen := createMenuScreen(&gameFont, mainForeground, accentBlue, &menuItems)

	// Set up parent screens
	creditScreen.SetParent(menuScreen)
	gameScreen.SetParent(menuScreen)
	settingsScreen.SetParent(menuScreen)

	var currentScreen ui.Screen = menuScreen

	for !rl.WindowShouldClose() {
		currentScreen = currentScreen.HandleInput()

		if _, ok := currentScreen.(*ui.QuitScreen); ok {
			break
		}

		rl.BeginDrawing()
		rl.ClearBackground(mainBackground)
		currentScreen.Draw()
		rl.EndDrawing()
	}
}
