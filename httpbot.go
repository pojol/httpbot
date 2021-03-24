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
	"github.com/pojol/httpbot/report"
	"github.com/pojol/httpbot/timeline"
)

type ErrInfo struct {
	ID  string
	Err error
}

// Bot HTTP测试机器人
type Bot struct {
	id       string
	name     string
	Timeline timeline.Timeline

	// 通用的http client，如果http request 的地址有区别，可以通过自定义 card.GetClient() 接口，来改变指向
	client *http.Client

	// 用于生成测试阶段的报告
	rep *report.Report

	// 元数据，用于保存整个http bot 测试阶段的数据
	meta interface{}

	createTime int64

	stop bool
}

// New new http test bot
func New(name string, client *http.Client, meta interface{}) *Bot {

	return &Bot{
		id:         uuid.New().String(),
		name:       name,
		meta:       meta,
		rep:        report.NewReport(),
		createTime: time.Now().Unix(),
		client:     client,
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
func (bot *Bot) Run(doneCh chan string, errCh chan ErrInfo) {

	go func() {
		var err error

		for _, s := range bot.Timeline.GetSteps() {

			for _, c := range s.GetCards() {
				if bot.stop {
					return
				}

				err = bot.exec(c)
				if err != nil {
					errCh <- ErrInfo{
						ID:  bot.id,
						Err: fmt.Errorf("%v %w", c.GetName(), err),
					}
					return
				}
			}

		}

		//bot.rep.Print()
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
	return bot.name
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
