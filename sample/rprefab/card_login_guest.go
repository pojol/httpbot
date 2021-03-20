package rprefab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pojol/httpbot/assert"
	"github.com/pojol/httpbot/prefab"
)

// LoginGuestRes 游客登录返回
type LoginGuestRes struct {
	Token string `json:"token"`
}

// LoginGuestCard 游客登录
type LoginGuestCard struct {
	Base  *prefab.Card
	URL   string
	delay time.Duration
	md    *BotDat
}

// NewGuestLoginCard 生成账号创建预制
func NewGuestLoginCard(md *BotDat) *LoginGuestCard {
	return &LoginGuestCard{
		Base:  prefab.NewCardWithConfig(),
		URL:   Urls[LoginGuest],
		delay: time.Millisecond,
		md:    md,
	}
}

// GetName 获取卡片名
func (card *LoginGuestCard) GetName() string { return LoginGuest }

// GetURL 获取服务器地址
func (card *LoginGuestCard) GetURL() string { return card.URL }

// GetClient 获取 http.client
func (card *LoginGuestCard) GetClient() *http.Client { return nil }

// GetHeader get header
func (card *LoginGuestCard) GetHeader() map[string]string { return card.Base.Header }

// GetMethod get method
func (card *LoginGuestCard) GetMethod() string { return card.Base.Method }

// SetDelay 设置卡片之间调用的延迟
func (card *LoginGuestCard) SetDelay(delay time.Duration) { card.delay = delay }

// GetDelay 获取卡片之间调用的延迟
func (card *LoginGuestCard) GetDelay() time.Duration { return card.delay }

// Enter 序列化传入消息体
func (card *LoginGuestCard) Enter() []byte {

	b := []byte{}

	card.Base.AddInjectAssert("token assert", func() error {
		return assert.NotEqual(card.md.Token, "")
	})

	return b
}

// Leave 反序列化返回消息
func (card *LoginGuestCard) Leave(res *http.Response) error {

	var err error
	var body []byte
	cres := LoginGuestRes{}

	errcode, _ := strconv.Atoi(res.Header["Errcode"][0])
	if errcode != 0 {
		err = fmt.Errorf("%v request err code = %v", card.GetURL(), errcode)
		goto EXT
	}

	body, _ = ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &cres)
	if err != nil {
		err = fmt.Errorf("%v json.Unmarshal err %v", card.GetURL(), err.Error())
		goto EXT
	}

	card.md.Token = cres.Token
	err = card.Base.Assert()

EXT:
	return err
}
