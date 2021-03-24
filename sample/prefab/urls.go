package prefab

// request url
const (
	LoginGuest = "LoginGuest"
)

// Urls 用于匹配覆盖率的请求列表
var Urls = map[string]string{
	LoginGuest: "http://127.0.0.1:14001/v1/login/guest",
}
