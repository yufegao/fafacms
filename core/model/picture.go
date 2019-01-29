package model

type Picture struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Url        string `json:"url"`
	CreateTime int    `json:"create_time,omitempty"`
	Status     int    `json:"status,omitempty xorm:"not null comment('1 normal，2 deleted') TINYINT(1)"`
	
	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}
