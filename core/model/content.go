package model

import (
	"errors"
	"github.com/hunterhug/fafacms/core/config"
)

// 内容表
type Content struct {
	Id                int    `json:"id" xorm:"bigint pk autoincr"`
	Seo               string `json:"seo" xorm:"index"`
	Title             string `json:"name" xorm:"varchar(200) notnull"`
	UserId            int    `json:"user_id" xorm:"bigint index"`                                                                         // 内容所属用户
	NodeId            int    `json:"node_id" xorm:"bigint index"`                                                                         // 节点ID
	Status            int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 ban, 3 rubbish，4 login delete') TINYINT(1) index"` // 0-1-2-3为正常
	Describe          string `json:"describe" xorm:"TEXT"`
	PreDescribe       string `json:"pre_describe" xorm:"TEXT"`                                            // 预览内容，临时保存，当修改后调用发布接口，会刷新到Describe，并且记录进历史表
	PreFlush          int    `json:"status" xorm:"not null comment('1 flush') TINYINT(1)"`                // 是否预览内容已经被刷新
	CloseComment      int    `json:"close_comment" xorm:"not null comment('0 close, 1 open') TINYINT(1)"` // 关闭评论开关，默认关闭
	Version           int    `json:"version"`                                                             // 发布了多少次版本
	CreateTime        int    `json:"create_time"`
	UpdateTime        int    `json:"update_time,omitempty"`
	ImagePath         string `json:"image_path" xorm:"varchar(700)"`
	Views             int    `json:"views"`                         // 被点击多少次，弱化
	SuggestUpdateTime int64  `json:"suggest_update_time,omitempty"` // 建议协程更新时间
	Good              int    `json:"good"`                          // 建议支持数量
	Bad               int    `json:"bad"`                           // 建议反对
	Ha                int    `json:"ha"`                            // 建议无所谓
	Password          string `json:"password,omitempty"`
	Aa                string `json:"aa,omitempty"`
	Ab                string `json:"ab,omitempty"`
	Ac                string `json:"ac,omitempty"`
	Ad                string `json:"ad,omitempty"`
}

// 内容历史表
type ContentHistory struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	ContentId  int    `json:"content_id" xorm:"bigint index"` // 内容ID
	Seo        string `json:"seo" xorm:"index"`
	Title      string `json:"name" xorm:"varchar(200) notnull"`
	UserId     int    `json:"user_id" xorm:"bigint index"` // 内容所属的用户ID
	NodeId     int    `json:"node_id" xorm:"bigint index"` // 内容所属的节点
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time"`
}

// 内容建议表，哪个用户对哪个内容进行了点评
type ContentSupport struct {
	Id            int `json:"id" xorm:"bigint pk autoincr"`
	UserId        int `json:"user_id" xorm:"bigint index"`         // 评论的客户ID
	ContentId     int `json:"content_id" xorm:"bigint index"`      // 内容ID
	ContentUserId int `json:"content_user_id" xorm:"bigint index"` // 内容拥有者的ID
	CreateTime    int `json:"create_time"`
	Suggest       int `json:"suggest" xorm:"not null comment('1 good, 0 Ha，2 bad') TINYINT(1) index"` // 评价
}

// 统计节点下的内容数量
func (c *Content) CountNumOfNode() (int64, int64, error) {
	if c.UserId == 0 || c.NodeId == 0 {
		return 0, 0, errors.New("where is empty")
	}

	// 不是逻辑删除的都找出来
	normalNum, err := config.FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("node_id=?", c.NodeId).And("status<?", 4).Count()
	if err != nil {
		return 0, 0, err
	}
	allNum, err := config.FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("node_id=?", c.NodeId).Count()
	if err != nil {
		return 0, 0, err
	}
	return allNum, normalNum, nil
}
