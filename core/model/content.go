package model

import (
	"errors"
	"github.com/hunterhug/fafacms/core/config"
	"time"
)

// 内容表
type Content struct {
	Id       int    `json:"id" xorm:"bigint pk autoincr"`
	Seo      string `json:"seo" xorm:"index"`
	Title    string `json:"title" xorm:"varchar(200) notnull"`
	UserId   int    `json:"user_id" xorm:"bigint index"` // 内容所属用户
	UserName string `json:"user_name" xorm:"index"`
	NodeId   int    `json:"node_id" xorm:"bigint index"` // 节点ID
	NodeSeo  string `json:"node_seo" xorm:"index"`       // 节点ID SEO
	// 0/1 ->3   2/3 -> 4
	Status       int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 ban, 3 rubbish，4 logic delete') TINYINT(1) index"` // 0-1-2-3为正常
	Top          int    `json:"top" xorm:"not null comment('0 normal, 1 top') TINYINT(1) index"`                                     // 置顶
	Describe     string `json:"describe" xorm:"TEXT"`
	PreDescribe  string `json:"pre_describe" xorm:"TEXT"`                                                           // 预览内容，临时保存，当修改后调用发布接口，会刷新到Describe，并且记录进历史表
	PreFlush     int    `json:"pre_flush" xorm:"not null comment('1 flush') TINYINT(1)"`                            // 是否预览内容已经被刷新
	CloseComment int    `json:"close_comment" xorm:"not null comment('0 close, 1 open, 2 direct open') TINYINT(1)"` // 关闭评论开关，默认关闭
	Version      int    `json:"version"`                                                                            // 0表示什么都没发布                                                      // 发布了多少次版本
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time,omitempty"`
	ImagePath    string `json:"image_path" xorm:"varchar(700)"`
	Views        int    `json:"views"` // 被点击多少次，弱化
	Password     string `json:"password,omitempty"`
	Aa           string `json:"aa,omitempty"`
	Ab           string `json:"ab,omitempty"`
	Ac           string `json:"ac,omitempty"`
	Ad           string `json:"ad,omitempty"`
}

var ContentSortName = []string{"=id", "-top", "-create_time", "-update_time", "-views", "=version", "+status", "=seo"}

// 内容历史表
type ContentHistory struct {
	Id         int    `json:"id" xorm:"bigint pk autoincr"`
	ContentId  int    `json:"content_id" xorm:"bigint index"` // 内容ID
	Seo        string `json:"seo" xorm:"index"`
	Title      string `json:"title" xorm:"varchar(200) notnull"`
	UserId     int    `json:"user_id" xorm:"bigint index"` // 内容所属的用户ID
	NodeId     int    `json:"node_id" xorm:"bigint index"` // 内容所属的节点
	Describe   string `json:"describe" xorm:"TEXT"`
	CreateTime int64  `json:"create_time"`
}

var ContentHistorySortName = []string{"-create_time"}

// 统计节点下的内容数量
func (c *Content) CountNumUnderNode() (int64, int64, error) {
	if c.UserId == 0 || c.NodeId == 0 {
		return 0, 0, errors.New("where is empty")
	}

	// 不是逻辑删除的都找出来
	normalNum, err := config.FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("node_id=?", c.NodeId).And("status<?", 4).Count()
	if err != nil {
		return 0, 0, err
	}
	allNum, err := config.FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("node_id=?", c.NodeId).Count()
	if err != nil {
		return 0, 0, err
	}
	return allNum, normalNum, nil
}

func (c *Content) CheckSeoValid() (bool, error) {
	// 用户ID和SEO必须同时存在
	if c.UserId == 0 || c.Seo == "" {
		return false, errors.New("where is empty")
	}

	// 常规统计
	num, err := config.FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("seo=?", c.Seo).Count()

	// 如果大于1表示存在
	if num >= 1 {
		return true, nil
	}
	return false, err
}

func (c *Content) Insert() (int64, error) {
	c.CreateTime = time.Now().Unix()
	return config.FafaRdb.InsertOne(c)
}

func (c *Content) Get() (bool, error) {
	if c.UserId == 0 || c.Id == 0 {
		return false, errors.New("where is empty")
	}

	// 逻辑删除的内容不能获取到
	return config.FafaRdb.Client.Where("status!=?", 4).Get(c)
}

func (c *Content) GetByAdmin() (bool, error) {
	if c.Id == 0 {
		return false, errors.New("where is empty")
	}

	if c.UserId != 0 {
		return c.Get()
	}

	return config.FafaRdb.Client.Get(c)
}

func (c *Content) GetByRaw() (bool, error) {
	return config.FafaRdb.Client.Get(c)
}

// 更新前都会调用 Get 接口
func (c *Content) Update() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	c.UpdateTime = time.Now().Unix()
	return config.FafaRdb.Client.MustCols("status", "close_comment", "pre_flush", "password", "top", "node_id", "node_seo").Omit("user_id").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

// 更新前都会调用 Get 接口
func (c *Content) UpdateView() {
	config.FafaRdb.Client.Where("id=?", c.Id).Incr("views").Cols("views").Update(c)
}

func (c *Content) UpdateDescribe() error {
	if c.UserId == 0 || c.Id == 0 {
		return errors.New("where is empty")
	}

	s := config.FafaRdb.Client.NewSession()
	if err := s.Begin(); err != nil {
		return err
	}

	defer s.Close()

	c.UpdateTime = time.Now().Unix()
	c.Version = c.Version + 1
	c.PreFlush = 1
	_, err := s.Cols("describe", "pre_flush", "update_time", "version").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
	if err != nil {
		s.Rollback()
		return err
	}

	ch := new(ContentHistory)
	ch.Seo = c.Seo
	ch.Describe = c.Describe
	ch.UserId = c.UserId
	ch.Title = c.Title
	ch.NodeId = c.NodeId
	ch.ContentId = c.Id
	ch.CreateTime = time.Now().Unix()
	_, err = s.InsertOne(ch)
	if err != nil {
		s.Rollback()
		return err
	}

	if err := s.Commit(); err != nil {
		s.Rollback()
		return err
	}
	return nil
}

func (c *Content) ResetDescribe() error {
	if c.UserId == 0 || c.Id == 0 {
		return errors.New("where is empty")
	}

	c.UpdateTime = time.Now().Unix()
	_, err := config.FafaRdb.Client.Cols("pre_describe", "update_time").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
	if err != nil {
		return err
	}

	return nil
}

// 0/1 ->3   2/3 -> 4
func (c *Content) UpdateStatusTo3() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	c.UpdateTime = time.Now().Unix()
	c.Status = 3
	return config.FafaRdb.Client.Cols("status", "update_time").Where("status<=?", 1).Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

func (c *Content) UpdateStatusTo3Reverse() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	c.UpdateTime = time.Now().Unix()
	c.Status = 0
	return config.FafaRdb.Client.Cols("status", "update_time").Where("status=?", 3).Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

func (c *Content) UpdateStatusTo4() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	c.UpdateTime = time.Now().Unix()
	c.Status = 4
	return config.FafaRdb.Client.Cols("status", "update_time").Where("status>=?", 2).Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

func (c *Content) UpdateStatus() (int64, error) {
	if c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	c.UpdateTime = time.Now().Unix()
	return config.FafaRdb.Client.Cols("status", "update_time").Where("id=?", c.Id).Update(c)
}

func (c *ContentHistory) GetByAdmin() (bool, error) {
	if c.Id == 0 {
		return false, errors.New("where is empty")
	}

	// 逻辑删除的内容不能获取到
	return config.FafaRdb.Client.Get(c)
}
