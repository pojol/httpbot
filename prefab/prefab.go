package prefab

import (
	"net/http"
	"time"
)

// ICard 逻辑卡片接口
type ICard interface {
	GetURL() string

	GetDelay() time.Duration
	SetDelay(delay time.Duration)

	GetHeader() map[string]string
	Marshal() []byte
	Unmarshal(res *http.Response) map[string]interface{}
}

// IMetaData 元数据
type IMetaData interface {
	Refresh(meta interface{})
}
