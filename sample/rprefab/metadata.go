package rprefab

type BotDat struct {
	Token string
	AccID string
}

func NewBotData() (*BotDat, error) {
	b := &BotDat{}
	return b, nil
}
