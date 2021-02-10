package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	bot "github.com/pojol/httpbot"
	"github.com/pojol/httpbot/sample/metadata"
	"github.com/pojol/httpbot/sample/steps"
)

var (
	help bool

	// target server addr
	target string

	// Strategy
	strategy string

	// robot number
	num int
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&target, "target", "http://123.207.198.57:2222", "set target server address")
	flag.StringVar(&strategy, "strategy", "normal", "set strategy")
	flag.IntVar(&num, "num", 1, "robot number")
}

func main() {
	initFlag()

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	//for i := 0; i < num; i++ {
	md := &metadata.BotMetaData{}
	bot := bot.New(bot.BotConfig{
		Addr: target,
	}, md)

	bot.Timeline.AddDelayStep(steps.NewAccLoginStep(md), time.Millisecond*100)
	bot.Timeline.AddLoopStep(steps.NewMailSendStep(md))

	bot.Run()
	//}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	bot.Report()
}
