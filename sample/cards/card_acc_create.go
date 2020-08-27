package cards

import (
	"encoding/json"
	"fmt"
)

// AccCreateCard 账号创建预制
type AccCreateCard struct {
	URL string
}

// NewAccCreateCard 生成账号创建预制
func NewAccCreateCard() *AccCreateCard {
	return &AccCreateCard{
		URL: "/acc.create",
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
	res := struct {
		Code int
		Msg  string
		Body []byte
	}{}

	err := json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println("card unmarshal err", err)
	}

	if res.Code != 200 {
		fmt.Println("card err", res.Code, res.Msg)
	}

	createRes := struct {
		Token string
	}{}
	err = json.Unmarshal(res.Body, &createRes)
	if err != nil {
		fmt.Println("card unmarshal err", err)
	}

	return map[string]interface{}{
		"acctoken": createRes.Token,
	}

}
