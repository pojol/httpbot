package prefab

import (
	"net/http"
	"reflect"
)

/* ICard
┌───────────┬──────────┬──────────┐
│ construct │  enter   │  leave   │
└───────────┴──────────┴──────────┘

#
* construct
* enter
* leave
*/
type ICard interface {
	GetName() string

	GetURL() string
	GetClient() *http.Client
	GetMethod() string
	GetHeader() map[string]string

	// InjectParm 注入参数
	// key : func () interface { return newkeyval }
	// 用于在 Enter 阶段进行参数注入
	InjectParm(key string, f func() interface{})

	// InjectAssert 注入断言
	// name : func () error { return assert.Equal(a , b) }
	// 用于在 Leave 阶段进行一些判定
	InjectAssert(name string, f func() error)

	Enter() []byte
	Leave(res *http.Response) error
}

type IInject interface {
	Inject(childptr interface{})
	Assert() error
}

type Card struct {
	parmInject   map[string]func() interface{}
	assertInject map[string]func() error

	method string
	header map[string]string
}

func NewCardWithConfig() *Card {
	cp := &Card{
		parmInject:   make(map[string]func() interface{}),
		assertInject: make(map[string]func() error),
		method:       "POST",
		header:       make(map[string]string),
	}

	cp.header["Content-type"] = "application/json"

	return cp
}

func (c *Card) InjectParm(key string, f func() interface{}) {
	c.parmInject[key] = f
}

func (c *Card) InjectAssert(name string, f func() error) {
	c.assertInject[name] = f
}

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
