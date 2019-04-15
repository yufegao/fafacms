package model

import (
	"errors"
	"github.com/hunterhug/fafacms/core/config"
	"time"
)

// Content --> ContentNode
type Content struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Seo        string `json:"seo" xorm:"index"`
	Title      string `json:"name" xorm:"varchar(200) notnull"`
	UserId     int    `json:"user_id" xorm:"index"` // who's
	NodeId     int    `json:"node_id" xorm:"index"`
	Status     int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 deleted') TINYINT(1) index"`
	Type       int    `json:"type" xorm:"not null comment('0 paper，1 photo') TINYINT(1) index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int    `json:"create_time"`
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
	UserId       int    `json:"user_id" xorm:"index"`
	Type         int    `json:"type" xorm:"not null comment('0 article,1 diary,2 photo') TINYINT(1) index"`
	Seo          string `json:"seo" xorm:"index"`
	Status       int    `json:"status" xorm:"not null comment('0 normal,1 hide,2 deleted') TINYINT(1) index"`
	Name         string `json:"name" xorm:"varchar(100) notnull"`
	Describe     string `json:"describe" xorm:"TEXT"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time,omitempty"`
	ImagePath    string `json:"image_path" xorm:"varchar(1000)"`
	ParentNodeId int    `json:"parent_node_id"`
	Level        int    `json:"level"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
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
	c, err := config.FafaRdb.Client.Table(n).Where("type=?", n.Type).Where("user_id=?", n.UserId).And("id=?", n.ParentNodeId).And("level=?", 0).Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (n *ContentNode) InsertOne() error {
	n.CreateTime = time.Now().Unix()
	_, err := config.FafaRdb.Insert(n)
	return err
}
