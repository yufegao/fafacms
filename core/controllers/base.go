package controllers

import "fmt"

var (
	LoginPermit     = 10000
	ParseJsonError  = 10001
	UploadFileError = 10002
)

var ErrorMap = map[int]string{
	LoginPermit:     "login permit",
	ParseJsonError:  "json parse err",
	UploadFileError: "upload file err",
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
