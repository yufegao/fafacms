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
	Name            string `json:"name" xorm:"varchar(100) notnull unique"` // id and name index
	NickName        string `json:"nick_name" xorm:"varchar(100) notnull"`
	Email           string `json:"email" xorm:"varchar(100) notnull unique"`
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
	Code            string `json:"code,omitempty"`         // forget password code
	CodeExpired     int64  `json:"code_expired,omitempty"` // forget password code expired

	// Future...
	Aa string `json:"aa,omitempty"`
	Ab string `json:"ab,omitempty"`
	Ac string `json:"ac,omitempty"`
	Ad string `json:"ad,omitempty"`
}

var UserSortName = []string{"=id", "=name", "-update_time", "-create_time", "-gender"}

func (u *User) Get() (err error) {
	var exist bool
	exist, err = config.FafaRdb.Client.Get(u)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("user not found")
	}
	return
}

func (u *User) Exist() (bool, error) {
	if u.Id == 0 && u.Name == "" && u.GroupId == 0 {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Count(u)

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (u *User) IsNameRepeat() (bool, error) {
	if u.Name == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Table(u).Where("name=?", u.Name).Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (u *User) IsEmailRepeat() (bool, error) {
	if u.Email == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Table(u).Where("email=?", u.Email).Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (u *User) InsertOne() error {
	u.CreateTime = time.Now().Unix()
	_, err := config.FafaRdb.Insert(u)
	return err
}

func (u *User) IsActivateCodeExist() (bool, error) {
	if u.ActivateMd5 == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Get(u)
	return c, err
}

func (u *User) UpdateStatus() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.UpdateTime = time.Now().Unix()
	_, err := config.FafaRdb.Client.Where("id=?", u.Id).Cols("status", "update_time").Update(u)
	return err
}

func (u *User) UpdateActivateCode() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.UpdateTime = time.Now().Unix()
	u.ActivateMd5 = util.GetGUID()
	u.ActivateExpired = time.Now().Add(48 * time.Hour).Unix()
	_, err := config.FafaRdb.Client.Where("id=?", u.Id).Cols("activate_md5", "activate_expired", "update_time").Update(u)
	return err
}

func (u *User) GetUserByEmail() (bool, error) {
	if u.Email == "" {
		return false, errors.New("where is empty")
	}
	c, err := config.FafaRdb.Client.Get(u)
	return c, err
}

func (u *User) UpdateCode() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.UpdateTime = time.Now().Unix()
	u.Code = util.GetGUID()[0:6]
	u.CodeExpired = time.Now().Unix() + 60
	_, err := config.FafaRdb.Client.Where("id=?", u.Id).Cols("code", "code_expired", "update_time").Update(u)
	return err
}

func (u *User) UpdatePassword() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.UpdateTime = time.Now().Unix()
	u.Code = ""
	u.CodeExpired = 0
	_, err := config.FafaRdb.Client.Where("id=?", u.Id).Cols("code", "code_expired", "update_time", "password").Update(u)
	return err
}

func (u *User) UpdateInfo() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}

	u.UpdateTime = time.Now().Unix()
	_, err := config.FafaRdb.Client.Where("id=?", u.Id).Update(u)
	return err
}
