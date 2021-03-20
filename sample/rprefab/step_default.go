package rprefab

import "github.com/pojol/httpbot/prefab"

// NewDefaultStep 创建默认的step
func NewDefaultStep(md *BotDat) *prefab.Step {
	step := prefab.NewStep()

	step.AddCard(NewGuestLoginCard(md))

	return step
}
