package model

type Comment struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	UserId     string `json:"user_id" xorm:"index"`
	ObjectId   int    `json:"object_id" xorm:"index"`
	ObjectUser int    `json:"object_user" xorm:"index"`
	CommentId  int    `json:"comment_id,omitempty"`
	Status     int    `json:"status" xorm:"not null comment('0 todo, 1 normal，2 deleted') TINYINT(1) index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time,omitempty"`
	Good       int    `json:"good"`
	Bad        int    `json:"bad"`
	Aa         string `json:"aa,omitempty"`
	Ab         string `json:"ab,omitempty"`
	Ac         string `json:"ac,omitempty"`
	Ad         string `json:"ad,omitempty"`
}
