package prefab

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// MailDat mail data
type MailDat struct {
	ID       string
	Title    string
	Received bool
}

// ItemDat item data
type ItemDat struct {
	ID  string
	Num int
}

// BotMetaData metadata struct
type BotMetaData struct {
	AccToken     string
	AccLoginTime int64
	AccBag       []ItemDat

	Mails []MailDat
}

// Refresh refresh data
func (md *BotMetaData) Refresh(meta interface{}) {

	err := mapstructure.Decode(meta, md)
	if err != nil {
		fmt.Println("refresh metadata err", err)
	}

}
