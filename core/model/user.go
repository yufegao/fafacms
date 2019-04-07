package model

import (
	"errors"
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
)

// User --> Group
// user can not delete
type User struct {
	Id              int    `json:"id" xorm:"bigint pk autoincr"`
	Name            string `json:"name" xorm:"varchar(100) notnull index"` // id and name index
	NickName        string `json:"nick_name" xorm:"varchar(100) notnull"`
	Email           string `json:"email" xorm:"varchar(100)"`
	WeChat          string `json:"wechat" xorm:"varchar(100)"`
	WeiBo           string `json:"weibo" xorm:"TEXT"`
	Github          string `json:"github" xorm:"TEXT"`
	QQ              string `json:"qq" xorm:"varchar(100)"`
	Password        string `json:"password,omitempty" xorm:"varchar(100)"`
	Gender          int    `json:"gender" xorm:"not null comment('0 unknow, 1 boy，2 girl') TINYINT(1)"`
	Describe        string `json:"describe" xorm:"TEXT"`
	HeadPhoto       string `json:"head_photo" xorm:"varchar(1000)"`
	HomeType        int    `json:"home_type" xorm:"not null comment('0 normal，2...') TINYINT(1)"` // looks what home page
	CreateTime      int    `json:"create_time"`
	UpdateTime      int    `json:"update_time,omitempty"`
	DeleteTime      int    `json:"delete_time,omitempty"`
	ActivateMd5     int    `json:"activate_md5,omitempty"`     // register and reset md5 to email
	ActivateExpired int    `json:"activate_expired,omitempty"` // md5 expired time
	Status          int    `json:"status" xorm:"not null comment('0 unactive, 1 normal，2 reset) TINYINT(1) index"`
	GroupId         int    `json:"group_id,omitempty" xorm:"index"`

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

func (m *User) Get(userId int) (err error) {
	var exist bool
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
