package prefab

import (
	"encoding/json"
	"fmt"
)

// AccLoginCard 账号登录预制
type AccLoginCard struct {
	URL string
	md  *BotMetaData
}

// NewAccLoginCard 生成账号登录预制
func NewAccLoginCard(md *BotMetaData) *AccLoginCard {
	return &AccLoginCard{
		URL: "/login",
		md:  md,
	}
}

// GetURL 获取服务器地址
func (card *AccLoginCard) GetURL() string { return card.URL }

// Marshal 序列化传入消息体
func (card *AccLoginCard) Marshal() []byte {

	fmt.Println("login", card.md.AccToken)

	b, err := json.Marshal(struct {
		Token string
	}{card.md.AccToken})
	if err != nil {
		fmt.Println("json.Marshal", err, card.GetURL())
		b = []byte{}
	}

	return b
}

// Unmarshal 反序列化返回消息
func (card *AccLoginCard) Unmarshal(data []byte) map[string]interface{} {
	return nil
}
