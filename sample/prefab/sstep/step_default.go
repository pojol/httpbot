package sstep

import (
	"github.com/pojol/httpbot/sample/prefab"
	"github.com/pojol/httpbot/sample/prefab/scard"
	"github.com/pojol/httpbot/timeline"
)

// NewDefaultStep 创建默认的step
func NewDefaultStep(md *prefab.BotDat) *timeline.Step {
	step := timeline.NewStep("NewDefaultStep")

	step.AddCard(scard.NewGuestLoginCard(md))

	return step
}
