package factory

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	bot "github.com/pojol/httpbot"
	"github.com/pojol/httpbot/internal"
)

// CreateBotFunc 创建机器人的工厂方法
type CreateBotFunc func(addr string, client *http.Client) *bot.Bot

// 机器人的运行模式
const (
	FactoryModeStatic   = "static"
	FactoryModeIncrease = "increase"
)

// 策略的选取模式
const (
	StrategyModeNormal   = "normal"
	StrategyModePriority = "priority"

	StrategyPickNormal = "normal"
	StrategyPickRandom = "random"
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

// StrategyInfo 策略信息
type StrategyInfo struct {
	Name string
	Urls []string
}

// BotFactory 机器人工厂
type BotFactory struct {
	normalFactory   map[string]CreateBotFunc
	priorityFactory []CreateBotFunc
	strategys       []StrategyInfo
	pickCursor      int

	parm Parm

	beginTime time.Time

	bots map[string]*bot.Bot

	client *http.Client

	report   Report
	urlMatch map[string]int

	lock sync.Mutex
	exit *internal.Switch
}

// Create 构建机器人
func Create(opts ...Option) (*BotFactory, error) {

	p := Parm{
		frameRate: time.Second * 1,
		mode:      FactoryModeStatic,
		lifeTime:  time.Hour,
		pickMode:  StrategyPickNormal,
	}

	for _, opt := range opts {
		opt(&p)
	}

	if len(p.addr) == 0 {
		return nil, errors.New("Undefine address")
	}

	f := &BotFactory{
		parm:          p,
		bots:          make(map[string]*bot.Bot),
		normalFactory: make(map[string]CreateBotFunc),
		exit:          internal.NewSwitch(),
		beginTime:     time.Now(),
		report: Report{
			urlMap: make(map[string]*urlDetail),
		},
		urlMatch: make(map[string]int),
	}

	for _, v := range p.matchUrl {
		f.urlMatch[v] = 0
	}

	if p.client == nil {
		client := &http.Client{}
		f.client = client
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
func (f *BotFactory) Append(strategy string, cbf CreateBotFunc, mode string) {

	appendflag := false

	if mode == StrategyModeNormal {
		if _, ok := f.normalFactory[strategy]; !ok {
			f.normalFactory[strategy] = cbf
			appendflag = true
		}
	} else if mode == StrategyModePriority {
		f.priorityFactory = append(f.priorityFactory, cbf)
		appendflag = true
	}

	if appendflag {
		f.strategys = append(f.strategys, StrategyInfo{
			Name: strategy,
		})
	}
}

// Report 输出报告
func (f *BotFactory) Report() {

	fmt.Println("+--------------------------------------------------------------------------------------------------------+")
	fmt.Printf("Req url%-53s Req count %-5s Average time %-5s Succ rate %-10s\n", "", "", "", "")

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

		fmt.Printf("%-60s %-15d %-18s %-10s %-5s\n", sk, v.reqNum, avg, succ, reqsize+" / "+ressize)
	}
	fmt.Println("+--------------------------------------------------------------------------------------------------------+")

	durations := int(time.Now().Sub(f.beginTime).Seconds())
	if durations <= 0 {
		durations = 1
	}

	qps := int(reqtotal / int64(durations))

	duration := strconv.Itoa(durations) + "s"
	fmt.Printf("robot : %d req count : %d duration : %s qps : %d errors : %d\n", f.report.botNum, f.report.reqNum, duration, qps, f.report.errNum)

	for k, v := range f.urlMatch {
		if v > 0 {
			fmt.Printf("%-60s match %v\n", k, v)
		} else {
			fmt.Printf("%-60s \033[1;31;40m%-10s\033[0m\n", k, "match 0")
		}
	}
}

func (f *BotFactory) pushReport(bot *bot.Bot) {
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

	var creater CreateBotFunc

	if len(f.strategys) <= 0 {
		return nil
	}

	// 先从优先池中选取
	if len(f.priorityFactory) != 0 {

		pf := f.priorityFactory[0]
		f.priorityFactory = f.priorityFactory[1:]

		creater = pf
	} else {
		strategys := []string{}
		for k := range f.normalFactory {
			strategys = append(strategys, k)
		}

		if len(strategys) == 0 {
			panic(errors.New("not normal strategys"))
		}

		var pickstrategy string
		if f.parm.pickMode == StrategyPickNormal {

			if f.pickCursor >= len(strategys) {
				f.pickCursor = 0
			}
			pickstrategy = strategys[f.pickCursor]
			f.pickCursor++

		} else {
			pickstrategy = strategys[rand.Intn(len(strategys))]
		}

		creater = f.normalFactory[pickstrategy]
	}

	bot := creater(f.parm.addr[rand.Intn(len(f.parm.addr))], f.client)
	return bot
}

// Run 运行
func (f *BotFactory) Run() error {

	if f.parm.tickCreateNum == 0 {
		f.parm.tickCreateNum = len(f.strategys)
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

func (f *BotFactory) static() {

	var wg sync.WaitGroup

	for i := 0; i < f.parm.tickCreateNum; i++ {
		bot := f.getRobot()
		f.bots[bot.ID()] = bot
		wg.Add(1)
		bot.Run(&wg)
	}

	wg.Wait()
	for _, bot := range f.bots {
		f.pushReport(bot)
	}
	f.Report()
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

				for i := 0; i < f.parm.tickCreateNum; i++ {

					bot := f.getRobot()
					bot.Run(nil)

					f.bots[bot.ID()] = bot
				}

			case <-f.exit.Done():
			}

			if f.exit.HasOpend() {
				for _, bot := range f.bots {
					bot.Close()
					f.pushReport(bot)
				}
				f.Report()
				return
			}
		}

	}()

}
