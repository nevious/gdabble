package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/utils"
)

type memberRole struct {
	member string
	role   string
}

type credit struct {
	credits        *[]memberRole
	parent         Screen
	font           *rl.Font
	fontColor      rl.Color
	highlightColor rl.Color
}

func (c *credit) Update() Screen {
	if rl.IsKeyPressed(rl.KeyBackspace) {
		return c.parent
	}
	return c
}

func (c *credit) Draw() {
	t := "Credits"
	var size int32 = 35
	offset := utils.GetCenterForText(t, size, *c.font)
	rl.DrawText(t, offset, 100, size, c.fontColor)
}

func (c *credit) SetParent(parentScreen Screen) Screen {
	c.parent = parentScreen
	return c
}

func NewCreditScreen(font *rl.Font, color, highlight rl.Color) Screen {
	return &credit{
		font:           font,
		fontColor:      color,
		highlightColor: highlight,
	}
}
