package prefab

// ICard 逻辑卡片接口
type ICard interface {
	GetURL() string

	Marshal() []byte
	Unmarshal(data []byte) map[string]interface{}
}

// IMetaData 元数据
type IMetaData interface {
	Refresh(meta interface{})
}
