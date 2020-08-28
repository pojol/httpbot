package gobot

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/pojol/gobot/botreport"
	"github.com/pojol/gobot/mapping"
	"github.com/pojol/gobot/prefab"
)

// BotConfig config
type BotConfig struct {
	Addr   string
	Report bool
}

// Bot http bot
type Bot struct {
	Timeline prefab.Timeline

	cfg BotConfig

	meta    prefab.IMetaData
	mapping *mapping.Mapping

	report *botreport.Report
}

// New new http test bot
func New(cfg BotConfig, meta prefab.IMetaData) *Bot {
	return &Bot{
		cfg:     cfg,
		meta:    meta,
		report:  botreport.NewReport(),
		mapping: mapping.NewMapping(),
	}
}

func (bot *Bot) exec(card prefab.ICard) {
	url := bot.cfg.Addr + card.GetURL()

	begin := time.Now().UnixNano()
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
		bot.report.SetInfo(card.GetURL(), false, int((time.Now().UnixNano()-begin)/1000/1000))
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		bot.report.SetInfo(card.GetURL(), true, int((time.Now().UnixNano()-begin)/1000/1000))
		m := card.Unmarshal(res)
		if m != nil {
			for k := range m {
				bot.mapping.Set(k, m[k])
			}
			bot.meta.Refresh(bot.mapping.GetAll())
		}

	} else {
		bot.report.SetInfo(card.GetURL(), false, int((time.Now().UnixNano()-begin)/1000/1000))
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
							time.Sleep(c.GetDelay())
						}

						time.Sleep(s.GetDelay())
					}
				}()

			} else {
				for _, c := range s.GetCards() {
					bot.exec(c)
					time.Sleep(c.GetDelay())
				}
			}

			time.Sleep(s.GetDelay())
		}

		if bot.cfg.Report {
			bot.Report()
		}
	}()

}

// Report print report
func (bot *Bot) Report() {
	bot.report.Print()
}
