package model

// Content --> ContentNode
type Content struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Title      string `json:"name" xorm:"varchar(200) notnull"`
	UserId     string `json:"user_id" xorm:"index"` // who's
	NodeId     int    `json:"node_id" xorm:"index"`
	Status     int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 deleted') TINYINT(1) index"`
	Type       int    `json:"type" xorm:"not null comment('0 paper，1 photo') TINYINT(1) index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time,omitempty"`
	DeleteTime int    `json:"delete_time,omitempty"`
	ImagePath  string `json:"image_path" xorm:"varchar(1000)"`
	Views      int    `json:"views"`
	Password   string `json:"password,omitempty"`
	Good       int    `json:"good"`
	Bad        int    `json:"bad"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

type ContentNode struct {
	Id           int    `json:"id" xorm:"bigint pk autoincr"`
	UserId       string `json:"user_id" xorm:"index"`
	Type         int    `json:"type" xorm:"not null comment('0 paper，1 photo') TINYINT(1) index"`
	Status       int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 deleted') TINYINT(1) index"`
	Name         string `json:"name" xorm:"varchar(100) notnull"`
	Describe     string `json:"describe" xorm:"TEXT"`
	CreateTime   int    `json:"create_time"`
	UpdateTime   int    `json:"update_time,omitempty"`
	ImagePath    string `json:"image_path" xorm:"varchar(1000)"`
	ParentNodeId int    `json:"parent_node_id"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}
