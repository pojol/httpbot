# httpbot

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/httpbot)](https://goreportcard.com/report/github.com/pojol/httpbot)
[![Doc Card](https://img.shields.io/badge/httpbot-doc-2ca5e0?style=flat&logo=appveyor)](https://pojol.gitbook.io/httpbot/)

[中文](README_CN.md)

<div align="center">
    <img src="https://i.postimg.cc/v86d0Vqv/image.png" alt="img" width="600">
</div>

### Feature
* **`Reusable`** HTTP request action, after defining the HTTP request, we can reuse this definition in any strategy (the request parameters can be changed by injection 
* We can **`arrange`** tests in different strategies and then conduct their own tests for specific scenes. 
* Provide factory methods, allowing users to use **`multiple drive models`** to test to achieve their respective tests in different scenarios. (C-> S self-test, integrated in the CI step for API test, pressure test, etc. ... 

### Component
* Prefab
  * **Metadata** Save the attribute variables used in the entire life cycle of the Bot
  * **Card** Used to wrap HTTP requests
* Arrange
  * **Timeline** Logic drives the timeline
  * **Step** Used to encapsulate different action items of the Bot. At this stage, you can inject `parameters` and `assertions` to control the behavior logic and detect right or wrong
  * **Strategy** Provide Bot creation method, and behavior choreography (mainly aggregate Step
* Driver
  * **Factory** Used for batch drives of bots



### Quick start
```go

bf, _ := factory.Create(
	factory.WithCreateNum(0),	// run all strategy
	factory.WithLifeTime(time.Minute),
	factory.WithRunMode(factory.FactoryModeStatic),
	factory.WithMatchUrl([]string{
		"/v1/login/guest",
		"/v1/base/account.info"
	}),
)
defer bf.Close()

bf.Append("default strategy", func(url string, client *http.Client) *httpbot.Bot {
	md, err := rprefab.NewBotData()
	if err != nil {
		panic(err)
	}

	bot := httpbot.New(md, 
		client, 
		httpbot.WithName("default bot"))

	defaultStep := prefab.NewStep()

	guestLoginCard := prefab.NewGuestLoginCard(md)
	guestLoginCard.Base.InjectAssert("token assert", func() error {
		return assert.NotEqual(md.Token, "")
	})
	defaultStep.AddCard(guestLoginCard)

	bot.Timeline.AddStep(step)

	return bot
})

bf.Run()

```



### Report
```shell
/v1/login/guest             Req count 1     Consume 26ms  Succ rate 1/1   0kb / 0kb

+------------------------------------------------------------------------------------------------+
Req url                                     Req count       Average time       Succ rate
/v1/login/guest             1               26ms               1/1        0kb / 0kb
+------------------------------------------------------------------------------------------------+
robot : 1 req count : 1 duration : 1s qps : 1 errors : 0

/v1/base/account.info             not match
coverage  1 / 2
```