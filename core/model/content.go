package model

import (
	"errors"
	"github.com/hunterhug/fafacms/core/config"
	"time"
)

type Content struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Seo        string `json:"seo" xorm:"index"`
	Title      string `json:"name" xorm:"varchar(200) notnull"`
	UserId     int    `json:"user_id" xorm:"index"`
	NodeId     int    `json:"node_id" xorm:"index"`
	Status     int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 deleted') TINYINT(1) index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time"`
	UpdateTime int    `json:"update_time,omitempty"`
	ImagePath  string `json:"image_path" xorm:"varchar(1000)"`
	Views      int    `json:"views"`
	Password   string `json:"password,omitempty"`
	Good       int    `json:"good"`
	Bad        int    `json:"bad"`
	Aa         string `json:"aa,omitempty"`
	Ab         string `json:"ab,omitempty"`
	Ac         string `json:"ac,omitempty"`
	Ad         string `json:"ad,omitempty"`
}

type ContentNode struct {
	Id           int    `json:"id" xorm:"bigint pk autoincr"`
	UserId       int    `json:"user_id" xorm:"index"`
	Seo          string `json:"seo" xorm:"index"`
	Status       int    `json:"status" xorm:"not null comment('0 normal,1 hide,2 deleted') TINYINT(1) index"`
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

var ContentNodeSortName = []string{"=id", "-update_time", "-create_time", "+status", "=seo"}

func (n *ContentNode) CheckSeoValid() (bool, error) {
	if n.UserId == 0 || n.Seo == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("seo=?", n.Seo).Count()
	if c >= 1 {
		return true, nil
	}
	return false, err
}

func (n *ContentNode) CheckParentValid() (bool, error) {
	if n.UserId == 0 || n.ParentNodeId == 0 {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("id=?", n.ParentNodeId).And("level=?", 0).Count()
	if c >= 1 {
		return true, nil
	}
	return false, err
}

func (n *ContentNode) CheckChildrenNum() (int, error) {
	if n.UserId == 0 || n.Id == 0 {
		return 0, errors.New("where is empty")
	}
	num, err := config.FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("parent_node_id=?", n.Id).Count()
	return int(num), err
}

func (n *ContentNode) InsertOne() error {
	n.CreateTime = time.Now().Unix()
	_, err := config.FafaRdb.Insert(n)
	return err
}

func (n *ContentNode) Get() (bool, error) {
	if n.Id == 0 && n.Seo == "" {
		return false, errors.New("where is empty")
	}
	return config.FafaRdb.Client.Get(n)
}

func (n *ContentNode) Update() error {
	if n.Id == 0 {
		return errors.New("where is empty")
	}
	n.UpdateTime = time.Now().Unix()
	_, err := config.FafaRdb.Client.Where("id=?", n.Id).Cols("seo", "level", "parent_node_id", "name", "describe", "update_time", "status", "image_path").Update(n)
	return err
}

func (n *ContentNode) Delete() error {
	if n.Id == 0 {
		return errors.New("where is empty")
	}
	_, err := config.FafaRdb.Client.Delete(n)
	return err
}

func (c *Content) CountNumOfNode() (int, error) {
	if c.UserId == 0 || c.NodeId == 0 {
		return 0, errors.New("where is empty")
	}
	num, err := config.FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("node_id=?", c.NodeId).Count()
	return int(num), err
}
