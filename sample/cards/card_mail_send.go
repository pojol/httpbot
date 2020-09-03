package cards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pojol/gobot/sample/metadata"
)

// MailSendCard mail send
type MailSendCard struct {
	URL   string
	delay time.Duration
	md    *metadata.BotMetaData
}

// MailDat mail dat
type MailDat struct {
	ID      string
	Title   string
	Content string
}

// NewMailSendCard new mail send card
func NewMailSendCard(md *metadata.BotMetaData) *MailSendCard {
	return &MailSendCard{
		URL:   "/mail.send",
		delay: time.Millisecond,
		md:    md,
	}
}

// GetURL 获取服务器地址
func (card *MailSendCard) GetURL() string { return card.URL }

// GetHeader get card header
func (card *MailSendCard) GetHeader() map[string]string { return nil }

// SetDelay 设置卡片之间调用的延迟
func (card *MailSendCard) SetDelay(delay time.Duration) { card.delay = delay }

// GetDelay 获取卡片之间调用的延迟
func (card *MailSendCard) GetDelay() time.Duration { return card.delay }

// Marshal 序列化传入消息体
func (card *MailSendCard) Marshal() []byte {

	b, err := json.Marshal(struct {
		Token string
		Mail  MailDat
	}{
		Token: card.md.AccToken,
		Mail: MailDat{
			Title:   "test" + strconv.Itoa(int(time.Now().Unix())),
			Content: "content",
		},
	})

	if err != nil {
		fmt.Println("json.Marshal", err, card.GetURL())
		b = []byte{}
	}

	return b
}

// Unmarshal 反序列化返回消息
func (card *MailSendCard) Unmarshal(res *http.Response) map[string]interface{} {

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

	mailRes := struct {
		Token string
		Mails []MailDat
	}{}
	err = json.Unmarshal(resDat.Body, &mailRes)
	if err != nil {
		fmt.Println("card unmarshal err", err)
	}

	return map[string]interface{}{
		"acctoken": mailRes.Token,
		"mails":    mailRes.Mails,
	}
}
