package steps

import (
	"gobot/prefab"
	"gobot/sample/cards"
	"time"
)

// NewAccLoginStep prefab setp
func NewAccLoginStep() *prefab.Step {
	step := prefab.NewStep()
	step.SetDelay(time.Millisecond * 100)

	step.AddCard(cards.NewAccCreateCard())

	return step
}
