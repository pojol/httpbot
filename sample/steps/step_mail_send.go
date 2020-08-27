package steps

import (
	"time"

	"github.com/pojol/gobot/prefab"
	"github.com/pojol/gobot/sample/cards"
	"github.com/pojol/gobot/sample/metadata"
)

// NewMailSendStep send mail step
func NewMailSendStep(md *metadata.BotMetaData) *prefab.Step {

	step := prefab.NewStep()
	step.SetLoop(time.Second)

	step.AddCard(cards.NewMailSendCard(md))
	step.AddCard(cards.NewAccInfoCard(md))

	return step
}
