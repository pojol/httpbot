package factory

import "testing"

func TestReport(t *testing.T) {
	f, err := Create(WithAddr([]string{"127.0.0.1"}), WithStrategy("default"))
	if err != nil {
		t.FailNow()
	}

	f.Report()
}
