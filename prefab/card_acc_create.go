package prefab

// AccCreateCard 账号创建预制
type AccCreateCard struct {
	URL string
}

// NewAccCreateCard 生成账号创建预制
func NewAccCreateCard() *AccCreateCard {
	return &AccCreateCard{
		URL: "/create",
	}
}

// GetURL 获取服务器地址
func (card *AccCreateCard) GetURL() string { return card.URL }

// Marshal 序列化传入消息体
func (card *AccCreateCard) Marshal() []byte {

	b := []byte{}

	return b
}

// Unmarshal 反序列化返回消息
func (card *AccCreateCard) Unmarshal(data []byte) map[string]interface{} {

	return map[string]interface{}{
		"acctoken": "xxx",
	}

}
