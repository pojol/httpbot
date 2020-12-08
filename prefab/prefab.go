package prefab

import (
	"net/http"
	"time"
)

// ICard 逻辑卡片接口
type ICard interface {
	GetURL() string
	GetMethod() string
	GetHeader() map[string]string

	GetDelay() time.Duration
	SetDelay(delay time.Duration)

	Marshal() []byte
	Unmarshal(res *http.Response)
}
