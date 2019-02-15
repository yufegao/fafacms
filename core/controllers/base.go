package controllers

import "fmt"

var (
	AuthPermit      = 10000
	ParseJsonError  = 10001
	UploadFileError = 10002
	LoginWrong      = 10003
	DBError         = 10004
	ParasError      = 10005
	NoLogin         = 10006

	Unknown = 99999
)

var ErrorMap = map[int]string{
	AuthPermit:      "auth permit",
	ParseJsonError:  "json parse err",
	UploadFileError: "upload file err",
	LoginWrong:      "username or password wrong",
	DBError:         "db operation err",
	ParasError:      "paras not right",
	NoLogin:         "no login",
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
