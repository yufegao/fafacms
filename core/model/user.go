package model

import (
	"errors"
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/util"
	"time"
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
	CreateTime      int64  `json:"create_time"`
	UpdateTime      int64  `json:"update_time,omitempty"`
	DeleteTime      int64  `json:"delete_time,omitempty"`
	ActivateMd5     string `json:"activate_md5,omitempty"`     // register and reset md5 to email
	ActivateExpired int64  `json:"activate_expired,omitempty"` // md5 expired time
	Status          int    `json:"status" xorm:"not null comment('0 unactive, 1 normal, 2 black') TINYINT(1) index"`
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

func (m *User) IsNameRepeat() (bool, error) {
	if m.Name == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Table(m).Where("name=?", m.Name).Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (m *User) IsEmailRepeat() (bool, error) {
	if m.Email == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Table(m).Where("email=?", m.Email).Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (m *User) InsertOne() error {
	_, err := config.FafaRdb.Insert(m)
	return err
}

func (m *User) IsCodeExist() (bool, error) {
	if m.ActivateMd5 == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Get(m)
	return c, err
}

func (m *User) UpdateStatus() error {
	if m.Id == 0 {
		return errors.New("where is empty")
	}
	m.UpdateTime = time.Now().Unix()
	_, err := config.FafaRdb.Client.Where("id=?", m.Id).Cols("status", "update_time").Update(m)
	return err
}

func (m *User) UpdateCode() error {
	if m.Id == 0 {
		return errors.New("where is empty")
	}
	m.UpdateTime = time.Now().Unix()
	m.ActivateMd5 = util.GetGUID()
	m.ActivateExpired = time.Now().Add(48 * time.Hour).Unix()
	_, err := config.FafaRdb.Client.Where("id=?", m.Id).Cols("activate_md5", "activate_expired", "update_time").Update(m)
	return err
}
