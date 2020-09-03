package steps

import (
	"github.com/pojol/gobot/prefab"
	"github.com/pojol/gobot/sample/cards"
)

// NewAccLoginStep prefab setp
func NewAccLoginStep() *prefab.Step {
	step := prefab.NewStep()

	step.AddCard(cards.NewAccCreateCard())

	return step
}
