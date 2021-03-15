package botreport

import (
	"errors"
	"testing"
)

func TestErrWrapper(t *testing.T) {
	r := NewReport()
	r.SetErr(errors.New("1s wrapping error"))
	r.SetErr(errors.New("2nd wrapping error"))

	r.Print()
}
