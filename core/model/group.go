package model

import (
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
)

type Group struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Name       string `json:"name,omitempty" xorm:"varchar(100) notnull"`
	Describe   string `json:"describe,omitempty" xorm:"TEXT`
	CreateTime int64  `json:"create_time,omitempty"`
	UpdateTime int64  `json:"update_time,omitempty"`
	ImagePath  string `json:"image_path" xorm:"TEXT`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

func (g *Group) Get(groupId int) (err error) {
	var exist bool
	g.Id = groupId
	exist, err = config.FafaRdb.Client.Get(g)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("group not found")
	}
	return
}
