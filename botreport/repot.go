package botreport

import (
	"errors"
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

// GetSuccRate get succ rate
func GetSuccRate(info []Info) string {
	c := len(info)
	succ := 0

	for _, v := range info {
		if v.State {
			succ++
		}
	}

	return strconv.Itoa(succ) + "/" + strconv.Itoa(c)
}

// GetAverageTime get info average time
func GetAverageTime(info []Info) (int, error) {

	c := len(info)
	if c == 0 {
		return 0, errors.New("not info")
	}

	sum := 0
	for _, v := range info {
		sum += v.Consume
	}

	return int(sum / c), nil
}

// Clear clear
func (r *Report) Clear() {

	for k := range r.Info {
		r.Info[k] = r.Info[k][:0]
	}
}

// Print print report
func (r *Report) Print() {

	for k := range r.Info {

		t, err := GetAverageTime(r.Info[k])
		if err != nil {
			continue
		}

		rate := GetSuccRate(r.Info[k])

		if t > 100 {
			fmt.Printf("%-30s Req count %-5d Average time \033[1;31;40m%-5s\033[0m Succ rate %-10s \n", k, len(r.Info[k]), strconv.Itoa(t)+"ms", rate)
		} else {

			fmt.Printf("%-30s Req count %-5d Average time %-5s Succ rate %-10s \n", k, len(r.Info[k]), strconv.Itoa(t)+"ms", rate)
		}

	}

	r.Clear()
}
