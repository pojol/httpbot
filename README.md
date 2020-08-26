# gobot
http test robot

---

### Feature
* prefab
* mapping
* timeline



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