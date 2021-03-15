package prefab

import (
	"net/http"
	"time"
)

type ICard interface {
	GetName() string

	GetDelay() time.Duration
	SetDelay(delay time.Duration)
}

type IClientCard interface {
	ICard

	GetURL() string
	GetClient() *http.Client
	GetMethod() string
	GetHeader() map[string]string

	Marshal() []byte
	Unmarshal(res *http.Response)
}

type IAssertCard interface {
	ICard

	// Do 执行断言
	Do() error
}
