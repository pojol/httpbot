package steps

import (
	"github.com/pojol/gobot/prefab"
	"github.com/pojol/gobot/sample/cards"
	"github.com/pojol/gobot/sample/metadata"
)

// NewAccLoginStep prefab setp
func NewAccLoginStep(md *metadata.BotMetaData) *prefab.Step {
	step := prefab.NewStep()

	step.AddCard(cards.NewAccCreateCard(md))

	return step
}
