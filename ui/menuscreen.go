package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/utils"
)

// Navigation Struct, mapping a label to a Screen Type
// -------------------------------------------------
type MenuItem struct {
	Label  string
	Screen Screen
}

// TODO: Move to dedicated file in case we want to quit from multiple places
// and continue to have things seperated by purpose.
// Quit Screen hack
// -------------------------------------------------
type QuitScreen struct{}

func (q *QuitScreen) SetParent(parent Screen) Screen { return q }
func (q *QuitScreen) Draw()                          {}
func (q *QuitScreen) Update() Screen                 { return q }

// Menu Structure
// -------------------------------------------------
type menu struct {
	menuItems      *[]MenuItem
	font           *rl.Font
	fontColor      rl.Color
	highlightColor rl.Color
	index          int
}

// menu is the root screen element, it does not have a parent
func (m *menu) SetParent(parent Screen) Screen { return m }

func (m *menu) Update() Screen {
	if rl.IsKeyPressed(rl.KeyEnter) {
		items := m.menuItems
		return (*items)[m.index].Screen
	}

	if rl.IsKeyPressed(rl.KeyJ) || rl.IsKeyPressed(rl.KeyDown) {
		m.index += 1
		if m.index >= len(*m.menuItems) {
			m.index = 0
		}
	} else if rl.IsKeyPressed(rl.KeyK) || rl.IsKeyPressed(rl.KeyUp) {
		m.index -= 1
		if m.index < 0 {
			m.index = len(*m.menuItems) - 1
		}
	}

	return m
}

func (m *menu) Draw() {
	t := "Main Menu"
	var size int32 = 25
	var vOffset int32 = 30 // vertical offset between choices
	offset := utils.GetCenterForText(t, size+10, *m.font)
	rl.DrawText(t, offset, 100, size+10, m.fontColor)

	start := int32(rl.GetScreenHeight() / 3)
	for i, item := range *m.menuItems {
		color := m.fontColor
		if m.index == i {
			color = m.highlightColor
		}
		// take offset with t, it's not centered perfectly but it gives a nice alignment
		// feature not a bug :blush:
		offset = utils.GetCenterForText(t, size, *m.font)
		rl.DrawText(
			item.Label,
			offset,
			start+vOffset*int32(i),
			size,
			color,
		)
	}
}

func NewMenuScreen(font *rl.Font, color, highlight rl.Color, items *[]MenuItem) Screen {
	return &menu{
		font:           font,
		fontColor:      color,
		highlightColor: highlight,
		index:          0,
		menuItems:      items,
	}
}
