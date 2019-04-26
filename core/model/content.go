package model

import (
	"errors"
	"github.com/hunterhug/fafacms/core/config"
)

// 内容表
type Content struct {
	Id          int    `json:"id" xorm:"bigint pk autoincr"`
	Seo         string `json:"seo" xorm:"index"`
	Title       string `json:"name" xorm:"varchar(200) notnull"`
	UserId      int    `json:"user_id" xorm:"index"`
	UserName    string `json:"user_name" xorm:"index"` // 冗余字段
	NodeId      int    `json:"node_id" xorm:"index"`
	Status      int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 deleted， 3 strong problem') TINYINT(1) index"` // 逻辑删除为2，删除后SEO要置空，3表示严重违规被禁
	Describe    string `json:"describe" xorm:"TEXT"`
	PreDescribe string `json:"pre_describe" xorm:"TEXT"`                             // 预览内容，临时保存，当修改后调用发布接口，会刷新到Describe，并且记录进历史表
	PreFlush    int    `json:"status" xorm:"not null comment('1 flush') TINYINT(1)"` // 是否预览内容已经被刷新
	CreateTime  int    `json:"create_time"`
	UpdateTime  int    `json:"update_time,omitempty"`
	ImagePath   string `json:"image_path" xorm:"varchar(1000)"`
	Views       int    `json:"views"`
	Password    string `json:"password,omitempty"`
	Aa          string `json:"aa,omitempty"`
	Ab          string `json:"ab,omitempty"`
	Ac          string `json:"ac,omitempty"`
	Ad          string `json:"ad,omitempty"`
}

// 内容历史表
type ContentHistory struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	ContentId  int    `json:"content_id" xorm:"bigint pk autoincr"`
	Seo        string `json:"seo" xorm:"index"`
	Title      string `json:"name" xorm:"varchar(200) notnull"`
	UserId     int    `json:"user_id" xorm:"index"`
	UserName   string `json:"user_name" xorm:"index"`
	NodeId     int    `json:"node_id" xorm:"index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time"`
}

// 内容支持表，哪个用户对哪个内容进行了点评
type ContentSupport struct {
	Id          int    `json:"id" xorm:"bigint pk autoincr"`
	UserId      int    `json:"user_id" xorm:"index"`
	UserName    string `json:"user_name" xorm:"index"`
	ContentId   int    `json:"content_id" xorm:"index"`
	ContentUser int    `json:"content_user" xorm:"index"`
	CreateTime  int    `json:"create_time"`
	Suggest     int    `json:"suggest" xorm:"not null comment('1 good, 0 Ha，2 bad') TINYINT(1) index"`
}

// 内容汇总表，内容的总评数量
type ContentCal struct {
	Id            int   `json:"id" xorm:"bigint pk autoincr"`
	ContentId     int   `json:"content_id" xorm:"index"`
	ContentUserId int   `json:"content_user_id" xorm:"index"`
	CreateTime    int   `json:"create_time"`
	UpdateTime    int64 `json:"update_time,omitempty"`
	Good          int   `json:"good"`
	Bad           int   `json:"bad"`
	Ha            int   `json:"ha"`
}

// 统计节点下的内容数量
// 被删除的内容不会被统计
func (c *Content) CountNumOfNode() (int, error) {
	if c.UserId == 0 || c.NodeId == 0 {
		return 0, errors.New("where is empty")
	}

	// 非删除状态下的内容
	num, err := config.FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("node_id=?", c.NodeId).And("status<=?", 2).Count()
	return int(num), err
}
