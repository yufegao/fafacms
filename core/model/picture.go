package model

type Picture struct {
	Id             int    `json:"id" xorm:"bigint pk autoincr"`
	Type           string `json:"type"`
	FileName       string `json:"file_name"`
	ReallyFileName string `json:"really_file_name"`
	Md             string `json:"md"`
	Url            string `json:"url"`
	Describe       string `json:"describe" xorm:"TEXT`
	CreateTime     int64  `json:"create_time,omitempty"`
	DeleteTime     int    `json:"delete_time,omitempty"`
	Status         int    `json:"status,omitempty xorm:"not null comment('1 normal，2 deleted') TINYINT(1)"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}
