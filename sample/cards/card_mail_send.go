package cards

import (
	"encoding/json"
	"fmt"
	"gobot/sample/metadata"
	"strconv"
	"time"
)

// MailSendCard mail send
type MailSendCard struct {
	URL string
	md  *metadata.BotMetaData
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
		URL: "/mail.send",
		md:  md,
	}
}

// GetURL 获取服务器地址
func (card *MailSendCard) GetURL() string { return card.URL }

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
func (card *MailSendCard) Unmarshal(data []byte) map[string]interface{} {
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

	mailRes := struct {
		Token string
		Mails []MailDat
	}{}
	err = json.Unmarshal(res.Body, &mailRes)
	if err != nil {
		fmt.Println("card unmarshal err", err)
	}

	return map[string]interface{}{
		"acctoken": mailRes.Token,
		"mails":    mailRes.Mails,
	}
}
