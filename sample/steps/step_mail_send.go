package steps

import (
	"github.com/pojol/httpbot/prefab"
	"github.com/pojol/httpbot/sample/cards"
	"github.com/pojol/httpbot/sample/metadata"
)

// NewMailSendStep send mail step
func NewMailSendStep(md *metadata.BotMetaData) *prefab.Step {

	step := prefab.NewStep()

	step.AddCard(cards.NewMailSendCard(md))
	step.AddCard(cards.NewAccInfoCard(md))

	return step
}
