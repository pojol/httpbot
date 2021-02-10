package steps

import (
	"github.com/pojol/httpbot/prefab"
	"github.com/pojol/httpbot/sample/cards"
	"github.com/pojol/httpbot/sample/metadata"
)

// NewAccLoginStep prefab setp
func NewAccLoginStep(md *metadata.BotMetaData) *prefab.Step {
	step := prefab.NewStep()

	step.AddCard(cards.NewAccCreateCard(md))

	return step
}
