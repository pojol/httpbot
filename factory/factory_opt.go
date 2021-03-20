package factory

import (
	"net/http"
	"time"
)

// Parm 机器人工厂可配置参数
type Parm struct {
	// lifeTime 工厂的生命周期
	//
	// 默认值 1h
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

	// pickMode 策略选取模式
	pickMode string

	// addr 目标网管地址
	addr []string

	// matchUrl 匹配路由列表
	matchUrl []string

	// client http client
	//
	// 如果没有调用 WithClient factory会创建一个默认的client
	client *http.Client
}

// Option consul discover config wrapper
type Option func(*Parm)

// WithAddr 目标服务器gate地址
func WithAddr(addr []string) Option {
	return func(c *Parm) {
		c.addr = addr
	}
}

// WithLifeTime 工厂的生命周期
func WithLifeTime(lifetime time.Duration) Option {
	return func(c *Parm) {
		c.lifeTime = lifetime
	}
}

// WithMode 运行模式
func WithMode(mode string) Option {
	return func(c *Parm) {
		c.mode = mode
	}
}

// WithFrameRate 工厂的运行帧频率
func WithFrameRate(framerate time.Duration) Option {
	return func(c *Parm) {
		c.frameRate = framerate
	}
}

// WithCreateNum 工厂每帧的机器人创建数
func WithCreateNum(num int) Option {
	return func(c *Parm) {
		c.tickCreateNum = num
	}
}

// WithMatchUrl 匹配url列表
func WithMatchUrl(urls []string) Option {
	return func(c *Parm) {
		c.matchUrl = urls
	}
}

// WithClient http client
func WithClient(client *http.Client) Option {
	return func(c *Parm) {
		c.client = client
	}
}

// WithStrategyPick 设置策略的选取模式
func WithStrategyPick(mode string) Option {
	return func(c *Parm) {
		c.pickMode = mode
	}
}
