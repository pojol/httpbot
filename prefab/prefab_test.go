package prefab

import (
	"net/http"
	"testing"

	"github.com/influxdata/influxdb/pkg/testing/assert"
	botassert "github.com/pojol/httpbot/assert"
)

type TestCard struct {
	BaseCard *Card

	Testparm1 string
	Testparm2 bool
}

func (tc *TestCard) GetName() string              { return "TestCard" }
func (tc *TestCard) GetURL() string               { return "" }
func (tc *TestCard) GetClient() *http.Client      { return nil }
func (tc *TestCard) GetMethod() string            { return tc.BaseCard.Method }
func (tc *TestCard) GetHeader() map[string]string { return tc.BaseCard.Header }

func (tc *TestCard) Enter() []byte {

	tc.BaseCard.Inject(tc)

	return []byte{}
}

func (tc *TestCard) Leave(res *http.Response) error {

	var err error
	err = tc.BaseCard.Assert()

	return err
}

func TestInjectParm(t *testing.T) {

	tc := &TestCard{
		BaseCard: NewCardWithConfig(),
	}

	tc.BaseCard.AddInjectParm("Testparm1", func() interface{} {
		return "newtestparm1value"
	})
	tc.BaseCard.AddInjectParm("Testparm2", func() interface{} {
		return true
	})

	tc.BaseCard.AddInjectAssert("assert testparm1", func() error {
		return botassert.Equal(tc.Testparm1, "newtestparm1value")
	})
	tc.BaseCard.AddInjectAssert("assert testparm2", func() error {
		return botassert.Equal(tc.Testparm2, true)
	})

	tc.Enter()
	err := tc.Leave(nil)
	assert.Equal(t, err, nil)
}
