package cards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// AccCreateCard 账号创建预制
type AccCreateCard struct {
	URL   string
	delay time.Duration
}

// NewAccCreateCard 生成账号创建预制
func NewAccCreateCard() *AccCreateCard {
	return &AccCreateCard{
		URL:   "/acc.create",
		delay: time.Millisecond,
	}
}

// GetURL 获取服务器地址
func (card *AccCreateCard) GetURL() string { return card.URL }

// GetHeader get card header
func (card *AccCreateCard) GetHeader() map[string]string { return nil }

// SetDelay 设置卡片之间调用的延迟
func (card *AccCreateCard) SetDelay(delay time.Duration) { card.delay = delay }

// GetDelay 获取卡片之间调用的延迟
func (card *AccCreateCard) GetDelay() time.Duration { return card.delay }

// Marshal 序列化传入消息体
func (card *AccCreateCard) Marshal() []byte {

	b := []byte{}

	return b
}

// Unmarshal 反序列化返回消息
func (card *AccCreateCard) Unmarshal(res *http.Response) map[string]interface{} {

	body, _ := ioutil.ReadAll(res.Body)

	resDat := struct {
		Code int
		Msg  string
		Body []byte
	}{}

	err := json.Unmarshal(body, &resDat)
	if err != nil {
		fmt.Println("card unmarshal err", err)
	}

	if resDat.Code != 200 {
		fmt.Println("card err", resDat.Code, resDat.Msg)
	}

	createRes := struct {
		Token string
	}{}
	err = json.Unmarshal(resDat.Body, &createRes)
	if err != nil {
		fmt.Println("card unmarshal err", err)
	}

	return map[string]interface{}{
		"acctoken": createRes.Token,
	}

}
