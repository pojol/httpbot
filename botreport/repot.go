package botreport

import (
	"fmt"
	"strconv"
)

// Info request report
type Info struct {
	State   bool
	Consume int
}

// Report robot info
type Report struct {
	Info map[string][]Info
}

// NewReport new report
func NewReport() *Report {
	return &Report{
		Info: make(map[string][]Info),
	}
}

// SetInfo set request state & consume
func (r *Report) SetInfo(url string, state bool, consume int) {

	r.Info[url] = append(r.Info[url], Info{
		State:   state,
		Consume: consume,
	})

}

func (r *Report) getSuccRate(url string) string {
	c := len(r.Info[url])
	succ := 0

	for _, v := range r.Info[url] {
		if v.State {
			succ++
		}
	}

	return strconv.Itoa(succ) + "/" + strconv.Itoa(c)
}

func (r *Report) getAverageTime(url string) int {
	c := len(r.Info[url])
	sum := 0
	for _, v := range r.Info[url] {
		sum += v.Consume
	}

	return int(sum / c)
}

// Print print report
func (r *Report) Print() {

	for k := range r.Info {
		t := r.getAverageTime(k)
		if t > 100 {
			fmt.Printf("%-30s Req count %-5d Average time \033[1;31;40m%-5d\033[0m ms Succ rate %-10s \n", k, len(r.Info[k]), t, r.getSuccRate(k))
		} else {
			fmt.Printf("%-30s Req count %-5d Average time %-5d ms Succ rate %-10s \n", k, len(r.Info[k]), t, r.getSuccRate(k))
		}
	}

}
