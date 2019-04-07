package model

import (
	"errors"
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
)

type Group struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	Name       string `json:"name" xorm:"varchar(100) notnull index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time,omitempty"`
	ImagePath  string `json:"image_path" xorm:"varchar(1000)"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

var GroupSortName = []string{"=id", "=name", "-create_time", "=update_time"}

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

func (g *Group) Exist() (bool, error) {
	if g.Id == 0 && g.Name == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Count(g)

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (g *Group) Delete() error {
	if g.Id == 0 && g.Name == "" {
		return errors.New("where is empty")
	}
	_, err := config.FafaRdb.Client.Delete(g)

	return err
}

func (g *Group) Take() (bool, error) {
	ok, err := g.Exist()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	_, err = config.FafaRdb.Client.Get(g)
	return true, err
}
