package report

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Info request report
type Info struct {
	State   bool
	Consume int
	ReqSize int64
	ResSize int64
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
func (r *Report) SetInfo(cardurl string, state bool, consume int, reqbyt int64, resbyt int64) {

	u, _ := url.Parse(cardurl)

	r.Info[u.Path] = append(r.Info[u.Path], Info{
		State:   state,
		Consume: consume,
		ReqSize: reqbyt,
		ResSize: resbyt,
	})

}

// GetSuccRate get succ rate
func getSuccRate(info []Info) string {
	c := len(info)
	succ := 0

	for _, v := range info {
		if v.State {
			succ++
		}
	}

	return strconv.Itoa(succ) + "/" + strconv.Itoa(c)
}

func getReqSize(info []Info) string {

	var s int64
	var ss string

	for _, v := range info {
		s += v.ReqSize
	}

	ks := int(s / 1024)
	if ks < 1024 {
		ss = strconv.Itoa(ks) + "kb"
	} else {
		ss = strconv.Itoa(int(ks/1024)) + "mb"
	}

	return ss
}

func getResSize(info []Info) string {
	var s int64
	var ss string

	for _, v := range info {
		s += v.ResSize
	}

	ks := int(s / 1024)
	if ks < 1024 {
		ss = strconv.Itoa(ks) + "kb"
	} else {
		ss = strconv.Itoa(int(ks/1024)) + "mb"
	}

	return ss
}

// GetAverageTime get info average time
func getAverageTime(info []Info) (int, error) {

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

		t, err := getAverageTime(r.Info[k])
		if err != nil {
			continue
		}

		rate := getSuccRate(r.Info[k])
		reqres := getReqSize(r.Info[k]) + " / " + getResSize(r.Info[k])

		if t > 100 {
			fmt.Printf("%-40s Req count %-5d Consume \033[1;31;40m%-5s\033[0m Succ rate %-5s %-5s\n", k, len(r.Info[k]), strconv.Itoa(t)+"ms", rate, reqres)
		} else {
			fmt.Printf("%-40s Req count %-5d Consume %-5s Succ rate %-5s %-5s\n", k, len(r.Info[k]), strconv.Itoa(t)+"ms", rate, reqres)
		}

	}

}
