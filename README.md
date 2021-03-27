# httpbot
一个基于线性时间驱动, 可进行逻辑编排的HTTP测试机器人框架

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/httpbot)](https://goreportcard.com/report/github.com/pojol/httpbot)
[![Doc Card](https://img.shields.io/badge/httpbot-doc-2ca5e0?style=flat&logo=appveyor)](https://pojol.gitbook.io/httpbot/)



#### Component
* Prefab
  * **Metadata** 元数据，用于保存Bot整个生命周期中使用到的属性变量
  * **Card** 用于包装HTTP请求（一次定义即可到处使用
* Arrange
  * **Timeline** Bot的逻辑驱动时间轴
  * **Step** 步骤条，用于封装Bot的不同行动项，在这个阶段可以注入`参数`以及`断言`，来控制行为逻辑和检测对错
  * **Strategy** 提供Bot的创建方法，以及行为编排（主要是聚合Step
* Driver
  * **Factory** 工厂，用于按指定的方式批量执行Bot



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