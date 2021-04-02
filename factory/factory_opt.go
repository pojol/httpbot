package factory

import (
	"net/http"
	"time"

	bot "github.com/pojol/httpbot"
)

// 机器人的运行模式
const (
	FactoryModeStatic   = "static"
	FactoryModeIncrease = "increase"
)

// 策略的选取模式
const (
	StrategyPickNormal = "normal"
	StrategyPickRandom = "random"
)

// CreateBotFunc 创建机器人的工厂方法
type CreateBotFunc func(fmd interface{}, client *http.Client) *bot.Bot

// Parm 机器人工厂可配置参数
type Parm struct {
	// lifeTime 工厂的生命周期
	//
	// 默认值 1分钟
	lifeTime time.Duration

	// frameRate 机器人工厂的运行频率，（每秒创建多少个机器人
	//
	// 默认值 1s
	frameRate time.Duration

	// tickCreateNum 机器人工厂每个频率创建的数量
	//
	// 默认值 1
	tickCreateNum int

	// mode 机器人工厂的运行模式
	//
	// FactoryModeStatic 静态模式，这种模式将只会执行第一帧，通常用于一次性运行若干机器人
	//
	// FactoryModeIncrease 累增模式，这种模式下会按频率不断创建机器人，并在生命周期到时销毁改机器人
	//
	// 默认值 FactoryModeStatic
	mode string

	// Interrupt 当card遇到err的时候是否中断整个程序 （默认为否
	Interrupt bool

	// pickMode 策略选取模式
	pickMode string

	// matchUrl 匹配路由列表
	matchUrl []string

	// client http client
	//
	// 如果没有调用 WithClient factory会创建一个默认的client
	client *http.Client

	// batchSize 批次大小（用于控制goroutine的并发数量（默认1024
	batchSize int
}

// Option consul discover config wrapper
type Option func(*Parm)

// WithLifeTime 定义工厂的生命周期 (默认为1分钟
func WithLifeTime(lifetime time.Duration) Option {
	return func(c *Parm) {
		c.lifeTime = lifetime
	}
}

// WithMode 运行模式，
func WithRunMode(mode string) Option {
	return func(c *Parm) {
		c.mode = mode
	}
}

// WithFrameRate 工厂的运行帧频率，用于 increase 模式下，每次批量创建bot的时间间隔（默认1秒
func WithFrameRate(framerate time.Duration) Option {
	return func(c *Parm) {
		c.frameRate = framerate
	}
}

// WithCreateNum 工厂每帧的机器人创建数， 当值为0时数量将会被初始化为 append strategy 的数量
func WithCreateNum(num int) Option {
	return func(c *Parm) {
		c.tickCreateNum = num
	}
}

// WithMatchUrl 匹配url列表，传入的match url 将会被用于计算是否被测试到和测试覆盖率。
func WithMatchUrl(urls []string) Option {
	return func(c *Parm) {
		c.matchUrl = urls
	}
}

// WithClient http client，使用自定义的 http.client
func WithClient(client *http.Client) Option {
	return func(c *Parm) {
		c.client = client
	}
}

// WithStrategyPick 设置策略的选取模式 (默认 normal 顺序执行
func WithStrategyPick(mode string) Option {
	return func(c *Parm) {
		c.pickMode = mode
	}
}

// WithInterrupt 设置遇到错误是否中断程序
func WithInterrupt(interrupt bool) Option {
	return func(c *Parm) {
		c.Interrupt = interrupt
	}
}
