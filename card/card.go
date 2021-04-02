package card

import (
	"net/http"
	"reflect"
	"time"
)

// ICard Card的抽象接口
// 这里主要用于 http 属性接口，以及http request 不同状态阶段的抽象。
type ICard interface {
	GetName() string

	GetURL() string
	GetClient() *http.Client
	GetMethod() string
	GetHeader() map[string]string

	GetDelay() time.Duration

	Enter() []byte
	Leave(res *http.Response) error
}

type IInject interface {
	// InjectParm 注入参数
	// key : func () interface { return newkeyval }
	// 用于在 Enter 阶段进行参数注入
	AddInjectParm(key string, f func() interface{})

	// InjectAssert 注入断言
	// name : func () error { return assert.Equal(a , b) }
	// 用于在 Leave 阶段进行一些判定
	AddInjectAssert(name string, f func() error)

	// Inject
	// 执行参数注入，一般在Enter阶段调用
	Inject(childptr interface{})

	// Assert
	// 执行断言判定， 一般在Leave阶段调用
	Assert() error
}

// Card 用于描述一些基础数据
type Card struct {
	parmInject   map[string]func() interface{}
	assertInject map[string]func() error

	Method string
	Header map[string]string
}

// NewCardWithConfig 创建一个新的card
func NewCardWithConfig() *Card {
	cp := &Card{
		parmInject:   make(map[string]func() interface{}),
		assertInject: make(map[string]func() error),
		Method:       "POST",
		Header:       make(map[string]string),
	}

	cp.Header["Content-type"] = "application/json"

	return cp
}

// AddInjectParm 添加一个参数注入
func (c *Card) AddInjectParm(key string, f func() interface{}) {
	c.parmInject[key] = f
}

// AddInjectAssert 添加一个断言注入
func (c *Card) AddInjectAssert(name string, f func() error) {
	c.assertInject[name] = f
}

// Inject 注入新的参数数据
func (c *Card) Inject(childptr interface{}) {
	if len(c.parmInject) == 0 {
		return
	}

	t := reflect.TypeOf(childptr).Elem()
	v := reflect.ValueOf(childptr).Elem()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() { // 检测是否可导出字段

			// 从注入表中开始进行检查
			for parm, injectf := range c.parmInject {
				if parm == t.Field(i).Name {

					val := injectf()
					// 判断类型是否适配
					if reflect.ValueOf(val).Type() == v.Field(i).Type() {
						v.Field(i).Set(reflect.ValueOf(val))
					}

				}
			}

		}
	}
}

// Assert 执行断言判定
func (c *Card) Assert() error {
	var err error
	for _, v := range c.assertInject {
		err = v()
		if err != nil {
			return err
		}
	}

	return nil
}
