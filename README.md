# gobot
http test robot

---

[![image.png](https://i.postimg.cc/FHTwqDq6/image.png)](https://postimg.cc/5XFPQqf5)


### Quick start
```go

    type BotMetaData struct {
        AccToken string
        Mails []Mail
    } 

    bot := gobot.New(BotConfig{
        ServerAddr : []string {"127.0.0.1"},
    }, &BotMetaData{})

    step := bot.Timeline.AddStep(time.Millisecond * 100)
    step.AddCard(prefab.NewAccLoginCard())

    bot.Run()

```