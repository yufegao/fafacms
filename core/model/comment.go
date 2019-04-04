package model

// Comment --> Content
type Comment struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	UserId     string `json:"user_id"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time,omitempty"`
	DeleteTime int    `json:"delete_time,omitempty"`
	Status     int    `json:"status,omitempty" xorm:"not null comment('1 normal，2 deleted') TINYINT(1)"`
	ContentId  int    `json:"content_id"`
	CommentId  int    `json:"comment_id,omitempty"`
	Good       int    `json:"good"`
	Bad        int    `json:"bad"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}
