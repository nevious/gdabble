package ui

type Screen interface {
	HandleInput() Screen // Return active/selected Screen
	Draw()               // Refresh on tick
	SetParent(parent Screen) Screen
}
