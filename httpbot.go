package httpbot

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pojol/httpbot/botreport"
	"github.com/pojol/httpbot/prefab"
)

// BotConfig config
type BotConfig struct {
	Addr   string
	Name   string
	Report bool
}

// Bot http bot
type Bot struct {
	id string

	Timeline prefab.Timeline

	cfg BotConfig

	meta interface{}

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
		createTime: time.Now().Unix(),
		client:     client,
	}
}

func (bot *Bot) exec(card prefab.ICard) error {
	bot.Lock()
	defer bot.Unlock()

	url := card.GetURL()
	var err error
	var res *http.Response
	var cheader map[string]string

	begin := time.Now().UnixNano()
	byt := card.Enter()
	reqsize := int64(len(byt))

	req, err := http.NewRequest(card.GetMethod(), url, bytes.NewBuffer(byt))
	if err != nil {
		goto EXT
	}

	cheader = card.GetHeader()
	if cheader != nil {
		for k, v := range cheader {
			req.Header.Set(k, v)
		}
	}

	if card.GetClient() != nil {
		res, err = card.GetClient().Do(req)
	} else {
		res, err = bot.client.Do(req)
	}

	if err != nil {
		bot.rep.SetErr(fmt.Errorf("client do err %v", err.Error()))
		goto EXT
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		err = card.Leave(res)
		if err == nil {
			bot.rep.SetInfo(card.GetURL(), true, int((time.Now().UnixNano()-begin)/1000/1000), reqsize, res.ContentLength)
		}
	} else {
		err = fmt.Errorf("http status %v url = %v err", res.Status, url)
	}
EXT:

	return err
}

// Run run bot
func (bot *Bot) Run(wg *sync.WaitGroup) {

	go func() {

		var err error

		for _, s := range bot.Timeline.GetSteps() {

			for _, c := range s.Step.GetCards() {
				if bot.stop {
					return
				}

				err = bot.exec(c)
				if err != nil {
					panic(fmt.Errorf("%v panic err %w", c.GetName(), err))
				}
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

// Name get bot name
func (bot *Bot) Name() string {
	return bot.cfg.Name
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
