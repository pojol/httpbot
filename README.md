# httpbot
一个基于线性时间驱动, 可编排的HTTP测试机器人框架

[![image.png](https://i.postimg.cc/3RbpyHvc/image.png)](https://postimg.cc/G8G9NVzF)

#### Feature
* 可复用,随意装配的http请求 (card
* 整个`Bot`生命周期可引用的`metadata`
* 可注入`参数`（主要用于Enter阶段)`断言`（用于Leave阶段做判定) 到card中
* 支持工厂模式，可批量创建不同模式，生命周期的`Bot`
* 格式化的报表输出

#### 
* Metadata 
    - 元数据; 用于保存在bot整个生命周期中使用到的属性变量，通常每个card都会持有md的引用。
* Card
    - 用于模拟一次http请求，包含三个阶段（构建，进入，离开）分别用于初始化http参数，参数注入&打包发送结构，解包发送结构&执行注入的断言函数。
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


```