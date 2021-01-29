package assert

import (
	"testing"

	"github.com/influxdata/influxdb/pkg/testing/assert"
)

type Obj struct {
	aa string
	bb int
	cc bool
}

func TestEqual(t *testing.T) {

	err := Equal("a", "aa")
	assert.NotEqual(t, err, nil)

	err = Equal("a", "a")
	assert.Equal(t, err, nil)

	err = Equal(1, 11)
	assert.NotEqual(t, err, nil)

	err = Equal(1, 1)
	assert.Equal(t, err, nil)

	err = Equal(Obj{
		aa: "aa",
		bb: 11,
		cc: false,
	}, Obj{
		aa: "aa",
		bb: 11,
		cc: true,
	})
	assert.NotEqual(t, err, nil)

	err = Equal(Obj{
		aa: "aa",
		bb: 11,
		cc: true,
	}, Obj{
		aa: "aa",
		bb: 11,
		cc: true,
	})
	assert.Equal(t, err, nil)
}
