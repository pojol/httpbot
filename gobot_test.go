package gobot

import (
	"fmt"
	"testing"
	"time"

	"github.com/pojol/gobot/sample/metadata"
	"github.com/pojol/gobot/sample/steps"
)

func TestBot(t *testing.T) {

	md := &metadata.BotMetaData{}

	bot := New(BotConfig{
		Addr: "http://123.207.198.57:2222",
	}, md)

	bot.Timeline.AddStep(steps.NewAccLoginStep(md))

	bot.Run()

	time.Sleep(time.Second)
	fmt.Println("token", md.AccToken)
}

func TestMapping(t *testing.T) {

	bot := New(BotConfig{}, &metadata.BotMetaData{})

	bot.mapping.Set("acctoken", "xxx")

	bot.mapping.Set("mails", []metadata.MailDat{
		{ID: "1", Title: "test1", Content: "content1"},
		{ID: "2", Title: "test2", Content: "content2"},
	})

	bot.mapping.Print()
}
