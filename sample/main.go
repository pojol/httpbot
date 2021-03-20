package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pojol/httpbot"
	"github.com/pojol/httpbot/factory"
	"github.com/pojol/httpbot/sample/rprefab"
)

func main() {

	targeturl := ""
	client := &http.Client{}

	var matchUrls []string
	for _, v := range rprefab.Urls {
		matchUrls = append(matchUrls, v)
	}

	bf, _ := factory.Create(
		factory.WithAddr([]string{targeturl}),
		factory.WithCreateNum(0),
		factory.WithClient(client),
		factory.WithMatchUrl(matchUrls),
	)
	defer bf.Close()

	bf.Append("default", func(url string, client *http.Client) *httpbot.Bot {
		md, err := rprefab.NewBotData()
		if err != nil {
			panic(err)
		}

		bot := httpbot.New(httpbot.BotConfig{
			Name:   "default bot",
			Addr:   url,
			Report: false,
		}, client, md)

		bot.Timeline.AddStep(rprefab.NewDefaultStep(md))

		return bot
	})

	bf.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
