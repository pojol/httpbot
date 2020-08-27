package steps

import (
	"gobot/prefab"
	"gobot/sample/cards"
	"gobot/sample/metadata"
	"time"
)

// NewMailSendStep send mail step
func NewMailSendStep(md *metadata.BotMetaData) *prefab.Step {

	step := prefab.NewStep()
	step.SetLoop(time.Second)

	step.AddCard(cards.NewMailSendCard(md))
	step.AddCard(cards.NewAccInfoCard(md))

	return step
}
