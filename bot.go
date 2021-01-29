package bot

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

	client *http.Client

	rep *botreport.Report
	sync.Mutex
}

// New new http test bot
func New(cfg BotConfig, client *http.Client, meta interface{}) *Bot {

	return &Bot{
		id:         uuid.New().String(),
		cfg:        cfg,
		meta:       meta,
		rep:        botreport.NewReport(),
		mapping:    mapping.NewMapping(),
		createTime: time.Now().Unix(),
		client:     client,
	}
}

func (bot *Bot) exec(card prefab.ICard) {
	bot.Lock()
	defer bot.Unlock()

	url := bot.cfg.Addr + card.GetURL()

	begin := time.Now().UnixNano()
	byt := card.Marshal()
	reqsize := int64(len(byt))

	req, err := http.NewRequest(card.GetMethod(), url, bytes.NewBuffer(byt))
	if err != nil {
		fmt.Println("http.NewRequest", err)
		return
	}

	cheader := card.GetHeader()
	if cheader != nil {
		for k, v := range cheader {
			req.Header.Set(k, v)
		}
	}

	res, err := bot.client.Do(req)
	if err != nil {
		bot.rep.SetInfo(card.GetURL(), false, int((time.Now().UnixNano()-begin)/1000/1000), reqsize, res.ContentLength)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		card.Unmarshal(res)

		bot.rep.SetInfo(card.GetURL(), true, int((time.Now().UnixNano()-begin)/1000/1000), reqsize, res.ContentLength)

	} else {
		bot.rep.SetInfo(card.GetURL(), false, int((time.Now().UnixNano()-begin)/1000/1000), reqsize, res.ContentLength)
	}

}

// Run run bot
func (bot *Bot) Run(wg *sync.WaitGroup) {

	go func() {

		for _, s := range bot.Timeline.GetSteps() {

			for _, c := range s.Step.GetCards() {
				if bot.stop {
					return
				}

				bot.exec(c)
				time.Sleep(time.Millisecond * 10)
			}

		}

		bot.rep.Print()
		bot.Close()
		if wg != nil {
			wg.Done()
		}
	}()

}

// ID get bot id
func (bot *Bot) ID() string {
	return bot.id
}

// GetReprotInfo 获取报告信息
func (bot *Bot) GetReprotInfo() map[string][]botreport.Info {
	bot.Mutex.Lock()
	defer bot.Mutex.Unlock()

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

// GetReport 获取报告指针
func (bot *Bot) GetReport() *botreport.Report {
	bot.Mutex.Lock()
	defer bot.Mutex.Unlock()

	return bot.rep
}

// Close close
func (bot *Bot) Close() {
	bot.stop = true
}

// Closed 机器人是否终止运行
func (bot *Bot) Closed() bool {
	return bot.stop
}
