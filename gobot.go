package gobot

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
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

	meta    interface{}
	mapping *mapping.Mapping

	report *botreport.Report
	sync.Mutex
}

// New new http test bot
func New(cfg BotConfig, meta interface{}) *Bot {
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
	cheader := card.GetHeader()
	if cheader != nil {
		for k, v := range cheader {
			req.Header.Set(k, v)
		}
	}

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
		card.Unmarshal(res)

	} else {
		bot.report.SetInfo(card.GetURL(), false, int((time.Now().UnixNano()-begin)/1000/1000))
		fmt.Println("http error status code", res.Status, "url", url)
	}

}

// Run run bot
func (bot *Bot) Run() {

	if bot.cfg.Report {
		bot.Report()
	}

	for _, s := range bot.Timeline.GetSteps() {

		if s.Loop {

			go func() {
				for {
					for _, c := range s.Step.GetCards() {
						bot.exec(c)
						//time.Sleep(c.GetDelay())
					}

					//time.Sleep(s.Dura)
				}
			}()

		} else {
			for _, c := range s.Step.GetCards() {
				bot.exec(c)
				//time.Sleep(c.GetDelay())
			}
		}

		//time.Sleep(s.Dura)
	}

}

// GetReport get report
func (bot *Bot) GetReport() *botreport.Report {
	bot.Lock()
	defer bot.Unlock()

	return bot.report
}

// Report print report
func (bot *Bot) Report() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				bot.report.Print()
			default:
			}
		}

	}()
}
