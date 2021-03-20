# httpbot
一个基于线性时间驱动, 可编排的HTTP测试机器人框架

[![image.png](https://i.postimg.cc/3RbpyHvc/image.png)](https://postimg.cc/G8G9NVzF)

#### Feature
* 可复用,随意装配的http请求 (card
* 整个`Bot`生命周期可引用的`metadata`
* 可注入`参数`（主要用于Enter阶段)`断言`（用于Leave阶段做判定) 到card中
* 支持工厂模式，可批量创建不同模式，生命周期的`Bot`
* 格式化的报表输出

#### Component
* Metadata 
    - 元数据; 用于保存在bot整个生命周期中使用到的属性变量，通常每个card都会持有md的引用。
* Card
    - 用于模拟一次http请求，包含三个阶段（构建，进入，离开）分别用于初始化http参数，参数注入&打包req结构，解包res结构&执行注入的断言函数。
* Timeline
    - 执行bot行为逻辑的时间轴
* Step
    - 时间轴上的步骤条，用于区分到不同的时间片上。 另外在step中还可以编排card的执行逻辑（包括注入参数等
* Strategy
    - 提供bot的创建方法，其中主要定义了bot的行为逻辑
* Factory
    - 工厂; 用于按指定的方式批量执行bot

### Quick start
```go

	bf, _ := factory.Create(
		factory.WithAddr([]string{targeturl}),
		factory.WithCreateNum(0),
		factory.WithClient(client),
	)
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

		bot.Timeline.AddStep(rprefab.NewDefaultStep(md))

		return bot
	})

	bf.Run()

```


#### report
```shell
http://127.0.0.1:14001/v1/login/guest                   请求数 1     耗时 26ms  成功率 1/1   0kb / 0kb

+--------------------------------------------------------------------------------------------------------+
Req url                                                 Req count       Average time       Succ rate
http://127.0.0.1:14001/v1/login/guest                   1               26ms               1/1        0kb / 0kb
+--------------------------------------------------------------------------------------------------------+
robot : 1 req count : 1 duration : 1s qps : 1 errors : 0

http://127.0.0.1:14001/v1/login/guest                   match 1
coverage  1 / 1
```