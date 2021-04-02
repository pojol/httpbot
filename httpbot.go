package httpbot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pojol/httpbot/card"
	"github.com/pojol/httpbot/internal"
	"github.com/pojol/httpbot/report"
	"github.com/pojol/httpbot/timeline"
)

type ErrInfo struct {
	ID  string
	Err error
}

// Bot HTTP测试机器人
type Bot struct {
	id string

	parm Parm

	client *http.Client

	Timeline timeline.Timeline

	// 用于生成测试阶段的报告
	rep *report.Report

	// 元数据，用于保存整个http bot 测试阶段的数据
	meta interface{}

	createTime int64

	stop bool
}

// New new http test bot
func New(meta interface{}, client *http.Client, opts ...Option) *Bot {

	p := Parm{
		PrintReprot: true,
	}
	for _, opt := range opts {
		opt(&p)
	}

	if client == nil {
		panic(fmt.Errorf("client is unavailable"))
	}

	return &Bot{
		id:         uuid.New().String(),
		parm:       p,
		client:     client,
		meta:       meta,
		rep:        report.NewReport(),
		createTime: time.Now().Unix(),
	}

}

func (bot *Bot) exec(c card.ICard) error {

	url := c.GetURL()
	var err error
	var res *http.Response
	var cheader map[string]string

	begin := time.Now().UnixNano()
	byt := c.Enter()
	reqsize := int64(len(byt))

	req, err := http.NewRequest(c.GetMethod(), url, bytes.NewBuffer(byt))
	if err != nil {
		goto EXT
	}

	cheader = c.GetHeader()
	if cheader != nil {
		for k, v := range cheader {
			req.Header.Set(k, v)
		}
	}

	if c.GetClient() != nil {
		res, err = c.GetClient().Do(req)
	} else {
		res, err = bot.client.Do(req)
	}

	if err != nil {
		goto EXT
	}
	defer res.Body.Close()
	req.Body.Close()

	if res.StatusCode == http.StatusOK {
		err = c.Leave(res)
		if err == nil {
			bot.rep.SetInfo(c.GetURL(), true, int((time.Now().UnixNano()-begin)/1000/1000), reqsize, res.ContentLength)
		}
	} else {
		io.Copy(ioutil.Discard, res.Body)
		err = fmt.Errorf("http status %v url = %v err", res.Status, url)
	}
EXT:

	return err
}

// Run run bot
func (bot *Bot) Run(sw *internal.Switch, doneCh chan string, errCh chan ErrInfo) {

	go func() {
		var err error

		for _, s := range bot.Timeline.GetSteps() {

			for _, c := range s.GetCards() {
				if bot.stop || sw.HasOpend() {
					return
				}

				err = bot.exec(c)
				if err != nil {
					errCh <- ErrInfo{
						ID:  bot.id,
						Err: fmt.Errorf("%v err -> %w", c.GetName(), err),
					}
					return
				}

				time.Sleep(c.GetDelay())
			}
		}

		if bot.parm.PrintReprot {
			bot.rep.Print()
		}

		bot.Close()
		doneCh <- bot.id
	}()

}

// ID get bot id
func (bot *Bot) ID() string {
	return bot.id
}

// Name get bot name
func (bot *Bot) Name() string {
	return bot.parm.Name
}

// GetReprotInfo 获取报告信息
func (bot *Bot) GetReprotInfo() map[string][]report.Info {
	nmap := make(map[string][]report.Info)
	for k, om := range bot.rep.Info {
		narr := []report.Info{}
		for _, info := range om {
			narr = append(narr, info)
		}

		nmap[k] = narr
	}

	return nmap
}

// GetReport 获取报告指针
func (bot *Bot) GetReport() *report.Report {
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
