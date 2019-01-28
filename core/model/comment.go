package model

type Comment struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	UserId     string `json:"user_id"`
	Describe   string `json:"describe" xorm:"TEXT`
	CreateTime int    `json:"create_time,omitempty"`
	DeleteTime int    `json:"delete_time,omitempty"`
	Status     int    `json:"status,omitempty xorm:"not null comment('1 normal，2 deleted') TINYINT(1)"`
	PaperId    int    `json:"paper_id"`
	CommentId  int    `json:"comment_id"`
	Good       int    `json:"good"`
	Bad        int    `json:"bad"`
}
