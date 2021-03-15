package factory

import "testing"

func TestReport(t *testing.T) {
	f, err := Create(WithAddr([]string{"127.0.0.1"}))
	if err != nil {
		t.FailNow()
	}

	f.Report()
}
