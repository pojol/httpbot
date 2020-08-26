# gobot
http test robot

---
[![image.png](https://i.postimg.cc/Z5XFtxKH/image.png)](https://postimg.cc/PCz81ZDv)


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