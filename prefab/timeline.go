package prefab

import (
	"time"
)

// TimelineStep timeline step
type TimelineStep struct {
	ID   string
	Step *Step
	Dura time.Duration
}

// Timeline Time-based step bar
type Timeline struct {
	steps []*TimelineStep
}

// AddStep add step in timeline
func (tl *Timeline) AddStep(step *Step) {
	tl.steps = append(tl.steps, &TimelineStep{
		Step: step,
		Dura: time.Millisecond,
	})
}

// AddDelayStep add need delay step
func (tl *Timeline) AddDelayStep(step *Step, dura time.Duration) {
	tl.steps = append(tl.steps, &TimelineStep{
		Step: step,
		Dura: dura,
	})
}

// SetStepDura set step dura
func (tl *Timeline) SetStepDura(id string, dura time.Duration) {
	for k, v := range tl.steps {
		if v.ID == id {
			tl.steps[k].Dura = dura
			break
		}
	}
}

// GetSteps get steps
func (tl *Timeline) GetSteps() []*TimelineStep {
	return tl.steps
}

// Step A step in the timeline
type Step struct {
	cards []ICard
}

// NewStep new step
func NewStep() *Step {
	step := &Step{}
	return step
}

// AddCard add prefab logic card
func (s *Step) AddCard(card ICard) {
	s.cards = append(s.cards, card)
}

// GetCards get card
func (s *Step) GetCards() []ICard {
	return s.cards
}
