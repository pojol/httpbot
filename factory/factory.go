package factory

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	bot "github.com/pojol/httpbot"
	"github.com/pojol/httpbot/internal"
	"github.com/pojol/httpbot/internal/color"
	"github.com/pojol/httpbot/internal/sizewg"
)

type urlDetail struct {
	reqNum int
	errNum int
	avgNum int64

	reqSize int64
	resSize int64
}

// 1. 信息的展示方式
// qps tps 流量使用情况 机器人并发数量 当前错误数
//
// 2. 单位时间内能够并发的机器人数量（性能
// 3. 机器人的运行机制

// Report 工厂的运行报告
/*
+-------------------------------------------------------------------+
| url           req count           average time          succ rate |
|                                                                   |
|/base/market.buy       1           10ms                  1/1       |
|/login/guest           1           9ms                   1/1       |
|/base/account.useItem  1           5ms                   1/1       |
+-------------------------------------------------------------------+
robot : 100, req : 100000, duration : 10s, errors : 0
*/
type Report struct {
	botNum int
	reqNum int
	errNum int

	urlMap map[string]*urlDetail
}

// StrategyInfo 策略信息结构
type StrategyInfo struct {
	Name string
	F    CreateBotFunc
}

// BotFactory 机器人工厂
type BotFactory struct {
	strategyLst []StrategyInfo
	pickCursor  int

	parm Parm

	colorer *color.Color

	beginTime time.Time

	client *http.Client

	bots map[string]*bot.Bot

	translateCh chan *bot.Bot
	doneCh      chan string
	errCh       chan bot.ErrInfo

	report   Report
	urlMatch map[string]int

	batch sizewg.SizeWaitGroup

	lock sync.Mutex
	exit *internal.Switch
}

// Create 构建机器人工厂
func Create(opts ...Option) (*BotFactory, error) {

	p := Parm{
		frameRate: time.Second * 1,
		mode:      FactoryModeStatic,
		lifeTime:  time.Minute,
		pickMode:  StrategyPickNormal,
		Interrupt: true,
		batchSize: 1024,
	}

	for _, opt := range opts {
		opt(&p)
	}

	f := &BotFactory{
		parm:      p,
		bots:      make(map[string]*bot.Bot),
		exit:      internal.NewSwitch(),
		beginTime: time.Now(),
		report: Report{
			urlMap: make(map[string]*urlDetail),
		},
		urlMatch:    make(map[string]int),
		translateCh: make(chan *bot.Bot),
		doneCh:      make(chan string),
		errCh:       make(chan bot.ErrInfo),
		colorer:     color.New(),
		batch:       sizewg.New(p.batchSize),
	}

	for _, v := range p.matchUrl {
		u, err := url.Parse(v)
		if err != nil {
			panic(err)
		}
		f.urlMatch[u.Path] = 0
	}

	if p.client == nil {
		f.client = &http.Client{}
	} else {
		f.client = p.client
	}

	return f, nil
}

// Close 关闭机器人工厂
func (f *BotFactory) Close() {
	f.exit.Open()
	//f.client.CloseIdleConnections()
}

// Append 添加机器人的创建策略
func (f *BotFactory) Append(strategy string, cbf CreateBotFunc) {

	f.strategyLst = append(f.strategyLst, StrategyInfo{
		Name: strategy,
		F:    cbf,
	})

}

// Report 输出报告
func (f *BotFactory) Report() {

	f.lock.Lock()
	defer f.lock.Unlock()

	fmt.Println("+--------------------------------------------------------------------------------------------------------+")
	fmt.Printf("Req url%-33s Req count %-5s Average time %-5s Succ rate %-10s\n", "", "", "", "")

	arr := []string{}
	for k := range f.report.urlMap {
		arr = append(arr, k)
	}
	sort.Strings(arr)

	var reqtotal int64

	for _, sk := range arr {
		v := f.report.urlMap[sk]
		avg := strconv.Itoa(int(v.avgNum/int64(v.reqNum))) + "ms"
		succ := strconv.Itoa(v.reqNum-v.errNum) + "/" + strconv.Itoa(v.reqNum)

		reqsize := strconv.Itoa(int(v.reqSize/1024)) + "kb"
		ressize := strconv.Itoa(int(v.resSize/1024)) + "kb"

		reqtotal += int64(v.reqNum)

		if _, ok := f.urlMatch[sk]; ok {
			f.urlMatch[sk]++
		}

		u, _ := url.Parse(sk)
		fmt.Printf("%-40s %-15d %-18s %-10s %-5s\n", u.Path, v.reqNum, avg, succ, reqsize+" / "+ressize)
	}
	fmt.Println("+--------------------------------------------------------------------------------------------------------+")

	durations := int(time.Now().Sub(f.beginTime).Seconds())
	if durations <= 0 {
		durations = 1
	}

	qps := int(reqtotal / int64(durations))

	duration := strconv.Itoa(durations) + "s"
	fmt.Printf("robot : %d req count : %d duration : %s qps : %d errors : %d\n", f.report.botNum, f.report.reqNum, duration, qps, f.report.errNum)

	if len(f.urlMatch) != 0 {
		coverage := 0
		for k, v := range f.urlMatch {
			if v > 0 {
				coverage++
				f.colorer.Printf("%-60s match\n", k)
			} else {
				f.colorer.Printf("%-60s %s\n", k, color.Red("not match"))
			}
		}

		if coverage == len(f.urlMatch) {
			f.colorer.Printf("coverage %v / %v\n", coverage, len(f.urlMatch))
		} else {
			f.colorer.Printf("coverage %v / %v\n", coverage, color.Red(len(f.urlMatch)))
		}
	}

}

func (f *BotFactory) pushReport(bot *bot.Bot) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.report.botNum++
	robotReport := bot.GetReprotInfo()

	for url, info := range robotReport {
		f.report.reqNum += len(info)

		if _, ok := f.report.urlMap[url]; !ok {
			f.report.urlMap[url] = &urlDetail{}
		}

		f.report.urlMap[url].reqNum += len(info)

		for _, v := range info {
			f.report.urlMap[url].avgNum += int64(v.Consume)
			f.report.urlMap[url].reqSize += v.ReqSize
			f.report.urlMap[url].resSize += v.ResSize
			if !v.State {
				f.report.errNum++
				f.report.urlMap[url].errNum++
			}
		}
	}
}

func (f *BotFactory) getRobot() *bot.Bot {

	if len(f.strategyLst) <= 0 {
		panic(errors.New("not strategys"))
	}

	var creator CreateBotFunc
	if f.parm.pickMode == StrategyPickNormal {
		if f.pickCursor >= len(f.strategyLst) {
			f.pickCursor = 0
		}
		creator = f.strategyLst[f.pickCursor].F
		f.pickCursor++

	} else {
		creator = f.strategyLst[rand.Intn(len(f.strategyLst))].F
	}

	bot := creator("", f.client)
	return bot
}

// Run 运行
func (f *BotFactory) Run() error {

	go f.router()

	if f.parm.tickCreateNum == 0 {
		f.parm.tickCreateNum = len(f.strategyLst)
	}

	if f.parm.mode == FactoryModeStatic {
		f.static()
	} else if f.parm.mode == FactoryModeIncrease {
		f.increase()
		time.AfterFunc(f.parm.lifeTime, func() {
			f.exit.Open()
		})
	}

	return nil
}

func (f *BotFactory) push(bot *bot.Bot) {
	f.batch.Add()

	f.bots[bot.ID()] = bot
}

func (f *BotFactory) pop(id string, err error) {
	f.batch.Done()

	if err != nil && f.parm.Interrupt {
		panic(err)
	}

	if _, ok := f.bots[id]; ok {

		if err == nil {
			f.pushReport(f.bots[id])
		} else {
			f.colorer.Printf("%v\n", color.Red(err.Error()))
		}

		delete(f.bots, id)

	}

	if len(f.bots) == 0 && f.parm.mode == FactoryModeStatic {
		f.exit.Open()
	}
}

func (f *BotFactory) router() {

	for {
		select {
		case bot := <-f.translateCh:
			f.push(bot)
			bot.Run(f.doneCh, f.errCh)
		case id := <-f.doneCh:
			f.pop(id, nil)
		case err := <-f.errCh:
			f.pop(err.ID, err.Err)
		case <-f.exit.Done():
			goto ext
		}
	}

ext:
	// clean
	// report
	f.Report()

	return
}

func (f *BotFactory) static() {

	for i := 0; i < f.parm.tickCreateNum; i++ {
		f.translateCh <- f.getRobot()
	}

	f.batch.Wait()

}

func (f *BotFactory) increase() {

	go func() {

		ticker := time.NewTicker(f.parm.frameRate)

		for {
			select {
			case <-ticker.C:

				if f.exit.HasOpend() {
					break
				}

				f.static()

			case <-f.exit.Done():
			}
		}

	}()

}
