package model

import (
	"errors"
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
)

type Resource struct {
	Id       int    `json:"id" xorm:"bigint pk autoincr"`
	Name     string `json:"name"`
	Url      string `json:"url" xorm:"varchar(1000) index"`
	Describe string `json:"describe" xorm:"TEXT"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

// extern: Group --> Resource
type GroupResource struct {
	Id         int `json:"id" xorm:"bigint pk autoincr"`
	GroupId    int `json:"group_id index"`
	ResourceId int `json:"resource_id index"`
}

func (r *Resource) Get() (err error) {
	var exist bool
	exist, err = config.FafaRdb.Client.Get(r)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("resource not found")
	}
	return
}

func (gr *GroupResource) Exist() (bool, error) {
	if gr.Id == 0 && gr.GroupId == 0 && gr.ResourceId == 0 {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Count(gr)

	if c >= 1 {
		return true, nil
	}

	return false, err
}
