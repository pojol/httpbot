package prefab

import "net/http"

// ICard 逻辑卡片接口
type ICard interface {
	GetURL() string

	Marshal() []byte
	Unmarshal(res *http.Response) map[string]interface{}
}

// IMetaData 元数据
type IMetaData interface {
	Refresh(meta interface{})
}
