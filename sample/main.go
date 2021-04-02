package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pojol/httpbot"
	"github.com/pojol/httpbot/factory"
	"github.com/pojol/httpbot/sample/prefab"
	"github.com/pojol/httpbot/sample/prefab/sstep"
)

var (
	help bool

	// target server addr
	target   []string
	backAddr string

	// strategy
	// default 用于测试临时流程，也可以自定义各种固有流程
	strategyParm string

	// robot number
	num int

	// lifttime 生命周期
	lifetime int

	// increase 增量
	increase bool
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	//target = append(target, "https://www.lunar-ring.com:14001")
	target = append(target, "https://www.lunar-ring.com:6443")

	flag.BoolVar(&increase, "increase", false, "incremental robot in every second")
	flag.IntVar(&lifetime, "lifetime", 60, "life time by second")
	flag.IntVar(&num, "num", 0, "robot number")
	flag.StringVar(&strategyParm, "strategy", "default", "robot strategy")
	flag.StringVar(&backAddr, "back", "", "back addr")
}

func main() {

	initFlag()

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	fmt.Println("bot num", num)
	fmt.Println("increase", increase)
	fmt.Println("lifetime", lifetime)

	rand.Seed(time.Now().UnixNano())

	mode := ""
	if increase {
		mode = factory.FactoryModeIncrease
	} else {
		mode = factory.FactoryModeStatic
	}

	var matchUrls []string
	for _, v := range prefab.Urls {
		matchUrls = append(matchUrls, v)
	}

	bf, _ := factory.Create(
		factory.WithCreateNum(num),
		factory.WithLifeTime(time.Duration(lifetime)*time.Second),
		factory.WithMatchUrl(matchUrls),
		factory.WithRunMode(mode),
	)
	defer bf.Close()

	bf.Append("default", func(fmd interface{}, client *http.Client) *httpbot.Bot {
		md, err := prefab.NewBotData()
		if err != nil {
			panic(err)
		}

		bot := httpbot.New(md, client,
			httpbot.WithName("default"))

		bot.Timeline.AddStep(sstep.NewDefaultStep(md))

		return bot
	})

	bf.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
