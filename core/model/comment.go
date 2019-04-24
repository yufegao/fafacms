package model

type Comment struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	UserId     string `json:"user_id" xorm:"index"`
	ObjectId   int    `json:"object_id" xorm:"index"`
	ObjectUser int    `json:"object_user" xorm:"index"`
	CommentId  int    `json:"comment_id,omitempty"`
	Status     int    `json:"status" xorm:"not null comment('1 normal, 0 hide，2 deleted') TINYINT(1) index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time,omitempty"`
	Aa         string `json:"aa,omitempty"`
	Ab         string `json:"ab,omitempty"`
	Ac         string `json:"ac,omitempty"`
	Ad         string `json:"ad,omitempty"`
}

type CommentSupport struct {
	Id           int    `json:"id" xorm:"bigint pk autoincr"`
	UserId       int    `json:"user_id" xorm:"index"`
	UserName     string `json:"user_name" xorm:"index"`
	ObjectId     int    `json:"object_id" xorm:"index"`
	ObjectUserId int    `json:"object_user_id" xorm:"index"`
	CommentId    int    `json:"comment_id" xorm:"index"`
	CreateTime   int    `json:"create_time"`
	Suggest      int    `json:"suggest" xorm:"not null comment('1 good, 0 Ha，2 bad') TINYINT(1) index"`
}

type CommentCal struct {
	Id           int   `json:"id" xorm:"bigint pk autoincr"`
	ObjectId     int   `json:"object_id" xorm:"index"`
	ObjectUserId int   `json:"object_user_id" xorm:"index"`
	CommentId    int   `json:"comment_id" xorm:"index"`
	CreateTime   int   `json:"create_time"`
	UpdateTime   int64 `json:"update_time,omitempty"`
	Good         int   `json:"good"`
	Bad          int   `json:"bad"`
	Ha           int   `json:"ha"`
}
