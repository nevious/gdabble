package ui

type Screen interface {
	Update() Screen // Return active/selected Screen
	Draw()          // Refresh on tick
	SetParent(parent Screen) Screen
}
