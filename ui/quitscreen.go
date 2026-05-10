package ui

type QuitScreen struct{}

func (q *QuitScreen) SetParent(parent Screen) Screen { return q }
func (q *QuitScreen) Render()                        {}
func (q *QuitScreen) Update() Screen                 { return q }
