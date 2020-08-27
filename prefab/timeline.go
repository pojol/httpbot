package prefab

import "time"

// Timeline Time-based step bar
type Timeline struct {
	steps []*Step
}

// AddStep add step in timeline
func (tl *Timeline) AddStep(step *Step) {
	tl.steps = append(tl.steps, step)
}

// GetSteps get steps
func (tl *Timeline) GetSteps() []*Step {
	return tl.steps
}

// Step A step in the timeline
type Step struct {
	cards []ICard
	loop  bool
	Dura  time.Duration
}

// NewStep new step
func NewStep() *Step {
	step := &Step{
		Dura: time.Millisecond * 1,
	}
	return step
}

// AddCard add prefab logic card
func (s *Step) AddCard(card ICard) {
	s.cards = append(s.cards, card)
}

// SetDelay set step delay
func (s *Step) SetDelay(delay time.Duration) {
	s.Dura = delay
}

// GetDelay get step delay
func (s *Step) GetDelay() time.Duration {
	return s.Dura
}

// SetLoop set step loop run
func (s *Step) SetLoop(delay time.Duration) {
	s.Dura = delay
	s.loop = true
}

// GetCards get card
func (s *Step) GetCards() []ICard {
	return s.cards
}

// Loop get loop setting
func (s *Step) Loop() bool {
	return s.loop
}
