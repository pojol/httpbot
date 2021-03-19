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
	ReqSize int64
	ResSize int64
}

// Report robot info
type Report struct {
	Info map[string][]Info
	err  error
}

// NewReport new report
func NewReport() *Report {
	return &Report{
		Info: make(map[string][]Info),
	}
}

// SetInfo set request state & consume
func (r *Report) SetInfo(url string, state bool, consume int, reqbyt int64, resbyt int64) {

	r.Info[url] = append(r.Info[url], Info{
		State:   state,
		Consume: consume,
		ReqSize: reqbyt,
		ResSize: resbyt,
	})

}

func (r *Report) SetErr(err error) {
	if r.err != nil {
		r.err = fmt.Errorf("%v =>\n%w", r.err.Error(), err)
	} else {
		r.err = fmt.Errorf("%w", err)
	}
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

	if r.err != nil {
		fmt.Println(r.err.Error())
		return
	}

	for k := range r.Info {

		t, err := getAverageTime(r.Info[k])
		if err != nil {
			continue
		}

		rate := getSuccRate(r.Info[k])
		reqres := getReqSize(r.Info[k]) + " / " + getResSize(r.Info[k])

		if t > 100 {
			fmt.Printf("%-60s 请求数 %-5d 耗时 \033[1;31;40m%-5s\033[0m 成功率 %-5s %-5s\n", k, len(r.Info[k]), strconv.Itoa(t)+"ms", rate, reqres)
		} else {
			fmt.Printf("%-60s 请求数 %-5d 耗时 %-5s 成功率 %-5s %-5s\n", k, len(r.Info[k]), strconv.Itoa(t)+"ms", rate, reqres)
		}

	}

}
