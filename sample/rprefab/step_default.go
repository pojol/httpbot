package rprefab

import "github.com/pojol/httpbot/prefab"

func NewDefaultStep(md *BotDat) *prefab.Step {
	step := prefab.NewStep()

	step.AddCard(NewGuestLoginCard(md))

	return step
}
