package gobot

import (
	"fmt"
	"gobot/prefab"
	"testing"
	"time"
)

func TestBot(t *testing.T) {

	var md prefab.BotMetaData

	bot := New(BotConfig{
		Addr: "127.0.0.1",
	}, &md)

	createStep := bot.Timeline.AddStep(time.Millisecond * 100)
	createStep.AddCard(prefab.NewAccCreateCard())

	loginStep := bot.Timeline.AddStep(time.Millisecond * 100)
	loginStep.AddCard(prefab.NewAccLoginCard(&md))

	bot.Run()

	time.Sleep(time.Second)
}

func TestMapping(t *testing.T) {

	var md prefab.BotMetaData
	bot := New(BotConfig{}, &md)

	bot.mapping.Set("acctoken", "xxx")
	bot.meta.Refresh(bot.mapping.GetAll())

	bot.mapping.Set("mails", []prefab.MailDat{
		{ID: "1", Title: "test1", Received: true},
		{ID: "2", Title: "test2", Received: true},
	})

	bot.meta.Refresh(bot.mapping.GetAll())

	bot.mapping.Set("acclogintime", time.Now().Unix())
	bot.meta.Refresh(bot.mapping.GetAll())

	bot.mapping.Set("accbag", []prefab.ItemDat{
		{ID: "item1", Num: 1},
		{ID: "item2", Num: 10}},
	)

	bot.meta.Refresh(bot.mapping.GetAll())

	fmt.Println("acc token", md.AccToken)
	fmt.Println("acc login time", md.AccLoginTime)
	fmt.Println("acc bag", md.AccBag)
	fmt.Println("mails", md.Mails)
}
