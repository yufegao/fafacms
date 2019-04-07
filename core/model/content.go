package model

// Content --> Node
type Content struct {
	Id     int    `json:"id" xorm:"bigint pk autoincr"`
	Title  string `json:"name,omitempty" xorm:"varchar(200) notnull"`
	UserId string `json:"user_id"`
	NodeId int    `json:"node_id"`
	Type   int    `json:"type,omitempty" xorm:"not null comment('1 paper，2 photo') TINYINT(1)"`
	Status int    `json:"status,omitempty" xorm:"not null comment('1 normal，2 deleted') TINYINT(1)"`

	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time,omitempty"`
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
	Name         string `json:"name,omitempty" xorm:"varchar(100) notnull"`
	Describe     string `json:"describe" xorm:"TEXT"`
	CreateTime   int    `json:"create_time,omitempty"`
	UpdateTime   int    `json:"update_time,omitempty"`
	ImagePath    string `json:"image_path" xorm:"varchar(1000)"`
	ParentNodeId int    `json:"parent_node_id"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}
