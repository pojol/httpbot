package prefab

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// MailDat mail data
type MailDat struct {
	ID      string
	Title   string
	Content string
}

// BotMetaData metadata struct
type BotMetaData struct {
	AccToken string
	Mails    []MailDat
}

// Refresh refresh data
func (md *BotMetaData) Refresh(meta interface{}) {

	err := mapstructure.Decode(meta, md)
	if err != nil {
		fmt.Println("refresh metadata err", err)
	}

}
