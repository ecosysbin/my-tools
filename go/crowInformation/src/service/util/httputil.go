package util

import (
	"github.com/astaxie/beego/httplib"
)

func GetResStr(url string) *httplib.BeegoHTTPRequest {
	return httplib.Get(url)
}
