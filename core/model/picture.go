package model

import (
	"errors"
	"github.com/hunterhug/fafacms/core/config"
)

type Picture struct {
	Id             int    `json:"id" xorm:"bigint pk autoincr"`
	Type           string `json:"type" xorm:"index"`
	Tag            string `json:"tag" xorm:"index"`
	UserId         int    `json:"user_id" xorm:"index"`
	FileName       string `json:"file_name"`
	ReallyFileName string `json:"really_file_name"`
	Md5            string `json:"md5" xorm:"index"`
	Url            string `json:"url" xorm:"varchar(1000) index"`
	Describe       string `json:"describe" xorm:"TEXT"`
	CreateTime     int64  `json:"create_time"`
	UpdateTime     int    `json:"update_time,omitempty"`
	Status         int    `json:"status" xorm:"not null comment('0 normal，1 hidebutcanuse') TINYINT(1)"`
	StoreType      int    `json:"store_type" xorm:"not null comment('0 local，1 oss') TINYINT(1)"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

func (p *Picture) Exist() (bool, error) {
	if p.Id == 0 && p.Url == "" && p.Md5 == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Count(p)

	if c >= 1 {
		return true, nil
	}

	return false, err
}
