package gobot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pojol/gobot/mapping"
	"github.com/pojol/gobot/prefab"
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
	Timeline prefab.Timeline

	cfg BotConfig

	meta    prefab.IMetaData
	mapping *mapping.Mapping

	// https://xxx:6443/xxx succ : 10, fail : 0, consume : 100ms
	Report map[string]Report
}

// New new http test bot
func New(cfg BotConfig, meta prefab.IMetaData) *Bot {
	return &Bot{
		cfg:     cfg,
		meta:    meta,
		mapping: mapping.NewMapping(),
	}
}

func (bot *Bot) exec(card prefab.ICard) {
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
		if m != nil {
			for k := range m {
				bot.mapping.Set(k, m[k])
			}
			bot.meta.Refresh(bot.mapping.GetAll())
		}

	} else {
		fmt.Println("http error status code", res.Status, "url", url)
	}

}

// Run run bot
func (bot *Bot) Run() {

	go func() {
		for _, s := range bot.Timeline.GetSteps() {

			if s.Loop() {

				go func() {
					for {
						for _, c := range s.GetCards() {
							bot.exec(c)
						}

						time.Sleep(s.GetDelay())
					}
				}()

			} else {
				for _, c := range s.GetCards() {
					bot.exec(c)
				}
			}

			time.Sleep(s.GetDelay())
		}
	}()

}
