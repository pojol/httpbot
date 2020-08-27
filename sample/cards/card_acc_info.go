package cards

import (
	"encoding/json"
	"fmt"

	"github.com/pojol/gobot/sample/metadata"
)

// AccInfoCard 账号登录预制
type AccInfoCard struct {
	URL string
	md  *metadata.BotMetaData
}

// NewAccInfoCard 查看账号信息
func NewAccInfoCard(md *metadata.BotMetaData) *AccInfoCard {
	return &AccInfoCard{
		URL: "/acc.info",
		md:  md,
	}
}

// GetURL 获取服务器地址
func (card *AccInfoCard) GetURL() string { return card.URL }

// Marshal 序列化传入消息体
func (card *AccInfoCard) Marshal() []byte {

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
func (card *AccInfoCard) Unmarshal(data []byte) map[string]interface{} {

	return nil
}
