package gobot

import (
	"gobot/prefab"
	"testing"
	"time"
)

func TestBot(t *testing.T) {

	md := &prefab.BotMetaData{}

	bot := New(BotConfig{
		Addr: "127.0.0.1",
	}, md)

	createStep := bot.Timeline.AddStep(time.Millisecond * 100)
	createStep.AddCard(prefab.NewAccCreateCard())

	loginStep := bot.Timeline.AddStep(time.Millisecond * 100)
	loginStep.AddCard(prefab.NewAccLoginCard(md))

	bot.Run()

	time.Sleep(time.Second)
}

func TestMapping(t *testing.T) {

	bot := New(BotConfig{}, &prefab.BotMetaData{})

	bot.mapping.Set("acctoken", "xxx")
	bot.meta.Refresh(bot.mapping.GetAll())

	bot.mapping.Set("mails", []prefab.MailDat{
		{ID: "1", Title: "test1", Content: "content1"},
		{ID: "2", Title: "test2", Content: "content2"},
	})
	bot.meta.Refresh(bot.mapping.GetAll())

	bot.mapping.Print()
}
