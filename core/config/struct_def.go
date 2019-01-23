package config

type Resp struct {
	Flag  bool        `json:"flag"`
	Error *ErrorResp  `json:"error"`
	Data  interface{} `json:"data"`
}

type ErrorResp struct {
	ErrorID  string `json:"id"`
	ErrorMsg string `json:"msg"`
}
