package model

type Log struct {
	Id           int    `json:"id" xorm:"bigint pk autoincr"`
	Cid          string `json:"cid"`
	Ip           string `json:"ip"`
	Url          string `json:"url"`
	LogTime      int64  `json:"log_time"`
	Ua           string `json:"ua"`
	UserId       int    `json:"user_id"`
	Flag         bool   `json:"flag"`
	In           string `json:"in"`
	Out          string `json:"out"`
	ErrorId      string `json:"error_id"`
	ErrorMessage string `json:"error_message"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}
