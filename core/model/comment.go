package model

// Comment --> Content
type Comment struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	UserId     string `json:"user_id" xorm:"index"`                                                          // who comment
	ObjectId   int    `json:"object_id" xorm:"index"`                                                        // comment which
	ObjectUser int    `json:"object_user" xorm:"index"`                                                      // comment who
	ObjectType int    `json:"object_type" xorm:"not null comment('0 content，1 photo') TINYINT(1) index"`     // comment type to decide object id
	CommentId  int    `json:"comment_id,omitempty"`                                                          // comment's comment
	Status     int    `json:"status" xorm:"not null comment('0 todo, 1 normal，2 deleted') TINYINT(1) index"` // delete is hide
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time"`
	DeleteTime int    `json:"delete_time,omitempty"`
	Good       int    `json:"good"`
	Bad        int    `json:"bad"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}
