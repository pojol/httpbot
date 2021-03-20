package rprefab

// BotDat bot的metadata
type BotDat struct {
	Token string
	AccID string
}

// NewBotData 创建bot metadata
func NewBotData() (*BotDat, error) {
	b := &BotDat{}
	return b, nil
}
