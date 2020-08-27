# gobot
http test robot

---
[![image.png](https://i.postimg.cc/Z5XFtxKH/image.png)](https://postimg.cc/PCz81ZDv)


### Quick start
```go

func NewMailSendStep(md *metadata.BotMetaData) *prefab.Step {

    step := prefab.NewStep()
    // set step run state
	step.SetLoop(time.Second)

    // add prefab logic card
	step.AddCard(cards.NewMailSendCard(md))
	step.AddCard(cards.NewAccInfoCard(md))

	return step
}

func main() {
    // define metadata
    type BotMetaData struct {
        AccToken string
        Mails []Mail
    } 

    // new bot
    bot := gobot.New(BotConfig{
        ServerAddr : []string {"127.0.0.1"},
    }, &BotMetaData{})

    // add prefab step
    bot.Timeline.AddStep(steps.NewMailSendStep())

    bot.Run()
}

```