package entity

/*
 * Define Actions or movement state an entity can be in
 */

type EntityMovementState int

const (
	StateIdle EntityMovementState = iota
	StateWalk
	StateAttack
)

var entityMovementStateNames = map[EntityMovementState]string{
	StateIdle:   "Idle",
	StateWalk:   "Walk",
	StateAttack: "Attack",
}

func (s EntityMovementState) String() string {
	return entityMovementStateNames[s]
}
