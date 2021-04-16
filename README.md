# httpbot
A http test robot framework that can be logically orchestrated

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/httpbot)](https://goreportcard.com/report/github.com/pojol/httpbot)
[![Doc Card](https://img.shields.io/badge/httpbot-doc-2ca5e0?style=flat&logo=appveyor)](https://pojol.gitbook.io/httpbot/)

<div align="center">
    <img src="https://i.postimg.cc/v86d0Vqv/image.png" alt="img" width="600">
</div>

#### Component
* Prefab
  * **Metadata** Save the attribute variables used in the entire life cycle of the Bot
  * **Card** Used to wrap HTTP requests
* Arrange
  * **Timeline** Logic drives the timeline
  * **Step** Used to encapsulate different action items of the Bot. At this stage, you can inject `parameters` and `assertions` to control the behavior logic and detect right or wrong
  * **Strategy** Provide Bot creation method, and behavior choreography (mainly aggregate Step
* Driver
  * **Factory** Used for batch drives of bots



#### Quick start
```go

bf, _ := factory.Create()
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

	defaultStep := prefab.NewStep()
	defaultStep.AddCard(prefab.NewGuestLoginCard(md))

	bot.Timeline.AddStep(step)

	return bot
})

bf.Run()

```



#### report
```shell
http://127.0.0.1:14001/v1/login/guest             Req count 1     Consume 26ms  Succ rate 1/1   0kb / 0kb

+--------------------------------------------------------------------------------------------------------+
Req url                                           Req count       Average time       Succ rate
http://127.0.0.1:14001/v1/login/guest             1               26ms               1/1        0kb / 0kb
+--------------------------------------------------------------------------------------------------------+
robot : 1 req count : 1 duration : 1s qps : 1 errors : 0

http://127.0.0.1:14001/v1/login/guest             match
coverage  1 / 1
```