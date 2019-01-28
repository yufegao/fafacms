package config

type Resp struct {
	Flag  bool        `json:"flag"`
	Error *ErrorResp  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type ErrorResp struct {
	ErrorID  string `json:"id"`
	ErrorMsg string `json:"msg"`
}
