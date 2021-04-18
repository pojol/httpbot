package timeline

import (
	"github.com/pojol/httpbot/card"
)

/*
	│                          Step
	├────────────────────────────►
	│
	│
	▼  Timeline
*/

// Step 步骤条, x轴，包含一组顺序执行的card
type Step struct {
	cards []card.ICard
	name  string
}

// Timeline 时间轴，按时间步进的y轴
type Timeline struct {
	steps []*Step
}

// AddStep 将一个步骤条添加到timeline中
func (tl *Timeline) AddStep(step *Step) {
	tl.steps = append(tl.steps, step)
}

// GetSteps 获取时间轴中的步骤条
func (tl *Timeline) GetCards() card.ICard {

	for sk, sv := range tl.steps {
		for ck, cv := range sv.cards {
			tl.steps[sk].cards = append(tl.steps[sk].cards[:ck], tl.steps[sk].cards[ck+1:]...)
			return cv
		}
	}

	return nil
}

// NewStep 创建一个新的步骤条
func NewStep(name string) *Step {
	step := &Step{
		name: name,
	}
	return step
}

// AddCard 添加一个预制的卡片到步骤条中
func (s *Step) AddCard(card card.ICard) {
	s.cards = append(s.cards, card)
}

func (s *Step) GetName() string {
	return s.name
}
