package controllers

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
)

var (
	LazyError       = 11111
	AuthPermit      = 10000
	ParseJsonError  = 10001
	UploadFileError = 10002
	LoginWrong      = 10003
	DBError         = 10004
	ParasError      = 10005
	NoLogin         = 10006
	DbNotFound      = 10007
	DbRepeat        = 10008
	DbHookIn        = 10009
	I500            = 99998
	Unknown         = 99999
)

var ErrorMap = map[int]string{
	DbNotFound:      "db not found",
	DbRepeat:        "db repeat data",
	DbHookIn:        "db hook in",
	I500:            "500 error",
	AuthPermit:      "auth permit",
	ParseJsonError:  "json parse err",
	UploadFileError: "upload file err",
	LoginWrong:      "username or password wrong",
	DBError:         "db operation err",
	ParasError:      "paras not right",
	NoLogin:         "no login",

	LazyError: "db not found or err",
}

func Error(code int, detail string) *ErrorResp {
	_, ok := ErrorMap[code]
	if !ok {
		code = Unknown
	}
	return &ErrorResp{
		ErrorID:  code,
		ErrorMsg: fmt.Sprintf("%s:%s", ErrorMap[code], detail),
	}
}

type Resp struct {
	Flag  bool        `json:"flag"`
	Cid   string      `json:"cid,omitempty"`
	Error *ErrorResp  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type ErrorResp struct {
	ErrorID  int    `json:"id"`
	ErrorMsg string `json:"msg"`
}

func (e ErrorResp) Error() string {
	return fmt.Sprintf("%d|%s", e.ErrorID, e.ErrorMsg)
}

type PageHelp struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Pages int `json:"total_pages"` // set by yourself
}

func (page *PageHelp) build(s *xorm.Session, sort []string, mapp map[string]string) {

	sortMapp := make(map[string]string)

	for _, v := range sort {
		if strings.HasPrefix(v, "+") {
			a := v[1:]
			if _, ok := mapp[a]; ok {
				sortMapp[a] = "+"
				s.Asc(a)
			}
		} else if strings.HasPrefix(v, "-") {
			a := v[1:]
			if _, ok := mapp[a]; ok {
				sortMapp[a] = "-"
				s.Desc(a)
			}
		}
	}

	for k, v := range mapp {
		if _, ok := sortMapp[k]; !ok {
			if v == "+" {
				s.Asc(k)
			} else if v == "-" {
				s.Desc(k)
			}
		}
	}

	if page.Page == 0 {
		page.Page = 1
	}

	if page.Limit <= 0 || page.Limit > 1000 {
		page.Limit = 20
	}
	s.Limit(page.Limit, (page.Page-1)*page.Limit)

}
