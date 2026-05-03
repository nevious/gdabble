package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/utils"
)

type settingsType int

const (
	ToggleSetting settingsType = iota
	SliderSetting
)

type Setting struct {
	name  string
	sType settingsType
}

type settingScreen struct {
	values    *[]Setting
	font      *rl.Font
	fontColor rl.Color
	highlight rl.Color
	parent    Screen
}

func (s *settingScreen) Draw() {
	t := "Settings"
	var size int32 = 35
	offset := utils.GetCenterForText(t, size, *s.font)
	rl.DrawText(t, offset, 100, size, s.fontColor)
}

func (s *settingScreen) SetParent(parent Screen) Screen {
	s.parent = parent
	return s
}

func (s *settingScreen) HandleInput() Screen {
	if rl.IsKeyPressed(rl.KeyBackspace) {
		return s.parent
	}
	return s
}

func NewSettingsScreen(values *[]Setting, font *rl.Font, color, highlight rl.Color) Screen {
	return &settingScreen{
		values:    values,
		font:      font,
		fontColor: color,
		highlight: highlight,
	}
}
