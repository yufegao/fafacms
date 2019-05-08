package controllers

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
)

// error code
var (
	LazyError         = 11111
	AuthPermit        = 10000
	ParseJsonError    = 10001
	UploadFileError   = 10002
	LoginWrong        = 10003
	DBError           = 10004
	ParasError        = 10005
	NoLogin           = 10006
	DbNotFound        = 10007
	DbRepeat          = 10008
	DbHookIn          = 10009
	EmailError        = 10010
	TimeNotReachError = 10011
	CodeWrong         = 10012
	I500              = 99998
	Unknown           = 99999
)

// error code message map
var ErrorMap = map[int]string{
	CodeWrong:         "code wrong",
	TimeNotReachError: "time not reach",
	EmailError:        "email error",
	DbNotFound:        "db not found",
	DbRepeat:          "db repeat data",
	DbHookIn:          "db hook in",
	I500:              "500 error",
	AuthPermit:        "auth permit",
	ParseJsonError:    "json parse err",
	UploadFileError:   "upload file err",
	LoginWrong:        "username or password wrong",
	DBError:           "db operation err",
	ParasError:        "paras not right",
	NoLogin:           "no login",

	LazyError: "db not found or err",
}

// common response
type Resp struct {
	Flag  bool        `json:"flag"`
	Cid   string      `json:"cid,omitempty"`
	Error *ErrorResp  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// inner error response
type ErrorResp struct {
	ErrorID  int    `json:"id"`
	ErrorMsg string `json:"msg"`
}

func (e ErrorResp) Error() string {
	return fmt.Sprintf("%d|%s", e.ErrorID, e.ErrorMsg)
}

func Error(code int, detail string) *ErrorResp {
	_, ok := ErrorMap[code]
	if !ok {
		code = Unknown
	}

	str := fmt.Sprintf("%s:%s", ErrorMap[code], detail)

	if detail == "" {
		str = fmt.Sprintf("%s", ErrorMap[code])
	}

	return &ErrorResp{
		ErrorID:  code,
		ErrorMsg: str,
	}
}

// list api page helper
type PageHelp struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Pages int `json:"total_pages"` // set by yourself outside
}

func (page *PageHelp) build(s *xorm.Session, sort []string, base []string) {
	Build(s, sort, base)

	if page.Page == 0 {
		page.Page = 1
	}

	if page.Limit <= 0 || page.Limit > 20 {
		page.Limit = 20
	}
	s.Limit(page.Limit, (page.Page-1)*page.Limit)
}

func Build(s *xorm.Session, sort []string, base []string) {
	for _, v := range base {
		a := v[1:]

		// if default use base sort field
		useBase := true
		for _, vv := range sort {
			b := vv[1:]

			// diy then change
			if a == b {
				useBase = false
				if strings.HasPrefix(vv, "+") {
					s.Asc(b)
				} else if strings.HasPrefix(vv, "-") {
					s.Desc(b)

				}
			}
		}

		if useBase {
			if strings.HasPrefix(v, "+") {
				s.Asc(a)
			} else if strings.HasPrefix(v, "-") {
				s.Desc(a)

			}
		}
	}
}
