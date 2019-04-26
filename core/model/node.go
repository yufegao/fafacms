package model

import (
	"errors"
	"github.com/hunterhug/fafacms/core/config"
	"time"
)

// 内容节点，最多两层
type ContentNode struct {
	Id           int    `json:"id" xorm:"bigint pk autoincr"`
	UserId       int    `json:"user_id" xorm:"index"`
	Seo          string `json:"seo" xorm:"index"`
	Status       int    `json:"status" xorm:"not null comment('0 normal,1 hide,2 deleted') TINYINT(1) index"` //  逻辑删除为2，SEO要置空
	Name         string `json:"name" xorm:"varchar(100) notnull"`
	Describe     string `json:"describe" xorm:"TEXT"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time,omitempty"`
	ImagePath    string `json:"image_path" xorm:"varchar(1000)"`
	ParentNodeId int    `json:"parent_node_id"`
	Level        int    `json:"level"`
	Aa           string `json:"aa,omitempty"`
	Ab           string `json:"ab,omitempty"`
	Ac           string `json:"ac,omitempty"`
	Ad           string `json:"ad,omitempty"`
}

// 内容节点排序专用，内容节点按更新时间降序，接着创建时间
var ContentNodeSortName = []string{"=id", "-update_time", "-create_time", "+status", "=seo"}

// 节点检查SEO的子路径是否有效
func (n *ContentNode) CheckSeoValid() (bool, error) {
	// 用户ID和SEO必须同时存在
	if n.UserId == 0 || n.Seo == "" {
		return false, errors.New("where is empty")
	}

	// 常规统计
	c, err := config.FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("seo=?", n.Seo).Count()

	// 如果大于1表示存在
	if c >= 1 {
		return true, nil
	}
	return false, err
}

// 节点检查 指定的父亲节点是否存在
func (n *ContentNode) CheckParentValid() (bool, error) {
	if n.UserId == 0 || n.ParentNodeId == 0 {
		return false, errors.New("where is empty")
	}

	// 只允许两层节点，Level必须为0
	// 被删除的节点不能统计
	c, err := config.FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("id=?", n.ParentNodeId).And("level=?", 0).And("status!=?", 2).Count()

	// 如果大于1表示存在
	if c >= 1 {
		return true, nil
	}
	return false, err
}

// 检查节点下的儿子节点数量
func (n *ContentNode) CheckChildrenNum() (int, error) {
	if n.UserId == 0 || n.Id == 0 {
		return 0, errors.New("where is empty")
	}
	num, err := config.FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("parent_node_id=?", n.Id).Count()
	return int(num), err
}

// 节点常规插入
func (n *ContentNode) InsertOne() error {
	n.CreateTime = time.Now().Unix()
	_, err := config.FafaRdb.Insert(n)
	return err
}

// 节点常规获取，ID和SEO必须存在一者
func (n *ContentNode) Get() (bool, error) {
	if n.Id == 0 && n.Seo == "" {
		return false, errors.New("where is empty")
	}
	return config.FafaRdb.Client.Get(n)
}

// 更新节点
func (n *ContentNode) Update() error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}
	n.UpdateTime = time.Now().Unix()
	_, err := config.FafaRdb.Client.Where("id=?", n.Id).And("user_id=?", n.UserId).Cols("seo", "level", "parent_node_id", "name", "describe", "update_time", "status", "image_path").Update(n)
	return err
}

// 逻辑删除，置空父亲和SEO
func (n *ContentNode) LogicDelete() error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}

	n.Status = 2
	n.UpdateTime = time.Now().Unix()
	n.Seo = ""
	n.ParentNodeId = 0
	n.Level = 0
	_, err := config.FafaRdb.Client.Where("id=?", n.Id).And("user_id=?", n.UserId).Cols("seo", "level", "parent_node_id", "update_time", "status").Update(n)
	return err
}
