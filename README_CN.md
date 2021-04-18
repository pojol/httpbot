# httpbot

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/httpbot)](https://goreportcard.com/report/github.com/pojol/httpbot)
[![Doc Card](https://img.shields.io/badge/httpbot-doc-2ca5e0?style=flat&logo=appveyor)](https://pojol.gitbook.io/httpbot/)

<div align="center">
    <img src="https://i.postimg.cc/v86d0Vqv/image.png" alt="img" width="600">
</div>

### 特性
* `可复用`的 HTTP 请求动作，在定义完 HTTP 请求之后，我们可以在任意的策略中复用这个定义（可以通过注入改变请求参数
* 逻辑`可编排`，我们可以将测试编排到各个不同的策略中，然后针对具体的场景进行各自的测试。
* 提供工厂方法，让用户可以采用`多种驱动模型`进行测试，以达到在不同的场景可以进行各自的测试。（C->S的自测，集成在CI步骤中进行API测试，压力测试 等等...

### 组件
* 预制阶段
  * **Metadata** 元数据，用于保存 Bot 在整个生命周期中使用的变量。
  * **Card** HTTP 请求的包装
* 编排阶段
  * **Timeline** 驱动逻辑执行顺序的时间轴
  * **Step** 用于封装Bot的不同行为。 在这一阶段，您可以注入 `参数` 和 `断言` 来控制行为逻辑并检测对与错 
  * **Strategy** 提供Bot的创建方法，以及行为编排（主要是整合Step
* 驱动阶段
  * **Factory** 用于批量创建Bot



### 快速开始
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



### 输出预览
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