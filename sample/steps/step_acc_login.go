package steps

import (
	"time"

	"github.com/pojol/gobot/prefab"
	"github.com/pojol/gobot/sample/cards"
)

// NewAccLoginStep prefab setp
func NewAccLoginStep() *prefab.Step {
	step := prefab.NewStep()
	step.SetDelay(time.Millisecond * 100)

	step.AddCard(cards.NewAccCreateCard())

	return step
}
