package gobot

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
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
	id string

	Timeline prefab.Timeline

	cfg BotConfig

	meta    interface{}
	mapping *mapping.Mapping

	stop       bool
	createTime int64

	rep *botreport.Report
	sync.Mutex
}

// New new http test bot
func New(cfg BotConfig, meta interface{}) *Bot {
	return &Bot{
		id:         uuid.New().String(),
		cfg:        cfg,
		meta:       meta,
		rep:        botreport.NewReport(),
		mapping:    mapping.NewMapping(),
		createTime: time.Now().Unix(),
	}
}

func (bot *Bot) exec(card prefab.ICard) {
	bot.Lock()
	defer bot.Unlock()

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

	//client := &http.Client{}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("client.Do", err)
		bot.rep.SetInfo(card.GetURL(), false, int((time.Now().UnixNano()-begin)/1000/1000))
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		card.Unmarshal(res)
		bot.rep.SetInfo(card.GetURL(), true, int((time.Now().UnixNano()-begin)/1000/1000))

	} else {
		bot.rep.SetInfo(card.GetURL(), false, int((time.Now().UnixNano()-begin)/1000/1000))
		fmt.Println("http error status code", res.Status, "url", url)
	}

}

// Run run bot
func (bot *Bot) Run() {

	if bot.cfg.Report {
		bot.report()
	}

	for _, s := range bot.Timeline.GetSteps() {

		if s.Loop {

			go func() {
				for {
					if bot.stop {
						break
					}

					for _, c := range s.Step.GetCards() {
						bot.exec(c)
						//time.Sleep(c.GetDelay())
					}

					time.Sleep(time.Millisecond * 100)
				}
			}()

		} else {
			for _, c := range s.Step.GetCards() {
				bot.exec(c)
				//time.Sleep(c.GetDelay())
			}
		}

		time.Sleep(time.Millisecond)
	}

}

// ID get bot id
func (bot *Bot) ID() string {
	return bot.id
}

// Close close
func (bot *Bot) Close() {
	bot.stop = true
}

// GetReportInfo get report
func (bot *Bot) GetReportInfo() map[string][]botreport.Info {
	bot.Lock()
	defer bot.Unlock()

	nmap := make(map[string][]botreport.Info)
	for k, om := range bot.rep.Info {
		narr := []botreport.Info{}
		for _, info := range om {
			narr = append(narr, info)
		}
		nmap[k] = narr
	}

	return nmap
}

// ClearReportInfo clear report info
func (bot *Bot) ClearReportInfo() {
	bot.Lock()
	defer bot.Unlock()

	bot.rep.Clear()
}

// Report print report
func (bot *Bot) report() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				bot.rep.Print()
			default:
			}
		}

	}()
}
