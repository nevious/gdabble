package ui

type Screen interface {
	// Update the screen state and return a screen
	Update() Screen
	// Render the Screen
	Render()
	// Set the Screens parent Screen
	SetParent(parent Screen) Screen
}
