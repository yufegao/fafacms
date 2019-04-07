package model

import (
	"errors"
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
)

// User --> Group
type User struct {
	Id          int    `json:"id" xorm:"bigint pk autoincr"`
	Name        string `json:"name" xorm:"varchar(100) notnull index"`
	NickName    string `json:"nick_name" xorm:"varchar(100) notnull"`
	Email       string `json:"email" xorm:"varchar(100)"`
	WeChat      string `json:"wechat" xorm:"varchar(100)"`
	WeiBo       string `json:"weibo" xorm:"TEXT"`
	QQ          string `json:"qq" xorm:"varchar(100)"`
	Password    string `json:"password,omitempty" xorm:"varchar(100)"`
	Gender      int    `json:"gender" xorm:"not null comment('0 unknow, 1 boy，2 girl') TINYINT(1)"`
	Describe    string `json:"describe" xorm:"TEXT"`
	ImagePath   string `json:"image_path" xorm:"varchar(1000)"`
	CreateTime  int    `json:"create_time"`
	UpdateTime  int    `json:"update_time,omitempty"`
	DeleteTime  int    `json:"delete_time,omitempty"`
	ActivateMd5 int    `json:"activate_md5,omitempty"`
	Status      int    `json:"status" xorm:"not null comment('0 unactive, 1 normal，2 reset, 3 deleted') TINYINT(1) index"`
	GroupId     int    `json:"group_id,omitempty" xorm:"index"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

func (m *User) Get(userId int) (err error) {
	var exist bool
	m.Status = 1
	m.Id = userId
	exist, err = config.FafaRdb.Client.Get(m)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("user not found")
	}
	return
}

func (m *User) Exist() (bool, error) {
	if m.Id == 0 && m.Name == "" && m.GroupId == 0 {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Count(m)

	if c >= 1 {
		return true, nil
	}

	return false, err
}
