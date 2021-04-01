package factory

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pojol/httpbot"
	"github.com/pojol/httpbot/card"
	"github.com/pojol/httpbot/sample/prefab"
	"github.com/pojol/httpbot/timeline"
)

type BCard struct {
	Base  *card.Card
	URL   string
	delay time.Duration
	md    *prefab.BotDat
}

func (card *BCard) GetName() string              { return "benchmark_card" }
func (card *BCard) GetURL() string               { return card.URL }
func (card *BCard) GetClient() *http.Client      { return nil }
func (card *BCard) GetHeader() map[string]string { return card.Base.Header }
func (card *BCard) GetMethod() string            { return card.Base.Method }
func (card *BCard) SetDelay(delay time.Duration) { card.delay = delay }
func (card *BCard) GetDelay() time.Duration      { return card.delay }
func (card *BCard) Enter() []byte                { return []byte{} }
func (card *BCard) Leave(res *http.Response) error {
	io.Copy(ioutil.Discard, res.Body)
	return nil
}

var bf *BotFactory

func TestMain(m *testing.M) {

	data := []byte("test")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write(data)
	}))
	defer srv.Close()

	bf, _ = Create(WithCreateNum(0), WithClient(&http.Client{}), WithRunMode(FactoryModeIncrease))
	bf.Append("benchmark_static", func(url string, client *http.Client) *httpbot.Bot {
		md, _ := prefab.NewBotData()
		bot := httpbot.New(md, client, httpbot.WithPrintReprot(false))

		step := timeline.NewStep("")
		step.AddCard(&BCard{
			Base:  card.NewCardWithConfig(),
			URL:   srv.URL,
			delay: time.Millisecond,
			md:    md,
		})
		bot.Timeline.AddStep(step)

		return bot
	})
	bf.parm.tickCreateNum = 1000
	go bf.router()

	m.Run()
	bf.Close()
}

// BenchmarkFactoryStatic-4   	   20532	     52667 ns/op
func BenchmarkFactoryStatic(b *testing.B) {

	//http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.static()
	}

}
