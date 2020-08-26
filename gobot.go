package gobot

import (
	"bytes"
	"fmt"
	"gobot/mapping"
	"io/ioutil"
	"net/http"
	"time"
)

// BotConfig config
type BotConfig struct {
	Addr string
}

// Report request report
type Report struct {
	URL     string
	Succ    int
	Fail    int
	Consume int
}

// Bot http bot
type Bot struct {
	Timeline Timeline

	cfg BotConfig

	meta    IMetaData
	mapping *mapping.Mapping

	// https://xxx:6443/xxx succ : 10, fail : 0, consume : 100ms
	Report map[string]Report
}

// Timeline bot timeline
type Timeline struct {
	steps []*Step
}

// Step A step in the timeline
type Step struct {
	cards []ICard
	Dura  time.Duration
}

// ICard 逻辑卡片接口
type ICard interface {
	GetURL() string

	Marshal() []byte
	Unmarshal(data []byte) map[string]interface{}
}

// IMetaData 元数据
type IMetaData interface {
	Refresh(meta interface{})
}

// New new http test bot
func New(cfg BotConfig, meta IMetaData) *Bot {
	return &Bot{
		cfg:     cfg,
		meta:    meta,
		mapping: mapping.NewMapping(),
	}
}

// AddStep add step in timeline
func (tl *Timeline) AddStep(dura time.Duration) *Step {

	step := &Step{
		Dura: dura,
	}
	tl.steps = append(tl.steps, step)

	return step
}

// AddCard add prefab logic card
func (step *Step) AddCard(card ICard) {
	step.cards = append(step.cards, card)
}

func (bot *Bot) exec(card ICard) {
	url := bot.cfg.Addr + card.GetURL()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(card.Marshal()))
	if err != nil {
		fmt.Println("http.NewRequest", err)
		return
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Println(url, "req succ")

		body, _ := ioutil.ReadAll(res.Body)

		m := card.Unmarshal(body)
		for k := range m {
			bot.mapping.Set(k, m[k])
		}

		bot.meta.Refresh(bot.mapping.GetAll())

	} else {
		fmt.Println("http error status code", res.Status, "url", url)
	}

}

// Run run bot
func (bot *Bot) Run() {

	go func() {
		for _, s := range bot.Timeline.steps {
			for _, c := range s.cards {
				bot.exec(c)
			}

			time.Sleep(s.Dura)
		}
	}()

}
